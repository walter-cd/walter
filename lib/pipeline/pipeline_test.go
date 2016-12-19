package pipeline

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/walter-cd/walter/lib/task"
)

func TestLoad(t *testing.T) {
	yaml := `
build:
  tasks:
    - name: command_stage_1
      command: echo "hello, world"
    - name: command_stage_2
      command: echo "hello, world, command_stage_2"
    - name: command_stage_3
      command: echo "hello, world, command_stage_3"
`
	p, err := Load([]byte(yaml))

	fmt.Printf("%v", p)
	if err != nil {
		t.Fatal(err)
	}

}

func TestSerialTasks(t *testing.T) {
	t1 := &task.Task{Name: "foo", Command: "echo foo"}
	t2 := &task.Task{Name: "bar", Command: "barbarbar"}
	t3 := &task.Task{Name: "baz", Command: "echo baz"}

	tasks := Tasks{t1, t2, t3}

	ctx, cancel := context.WithCancel(context.Background())

	p := &Pipeline{}
	p.runTasks(ctx, cancel, tasks, nil)

	if t1.Status != task.Succeeded {
		t.Fatal("t1 should have succeeded")
	}

	if t2.Status != task.Failed {
		t.Fatal("t2 should have failed")
	}

	if t3.Status != task.Skipped {
		t.Fatalf("t2 should have beed skipped")
	}
}

func TestSerialAndParallelTasks(t *testing.T) {
	p1 := &task.Task{Name: "p1", Command: "sleep 1"}
	p2 := &task.Task{Name: "p2", Command: "p2p2p2p2"}
	p3 := &task.Task{Name: "p3", Command: "sleep 1"}

	t1 := &task.Task{Name: "foo", Command: "echo foo"}
	t2 := &task.Task{Name: "bar", Parallel: Tasks{p1, p2, p3}}
	t3 := &task.Task{Name: "baz", Command: "echo baz"}

	tasks := Tasks{t1, t2, t3}
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pipeline{}
	p.runTasks(ctx, cancel, tasks, nil)

	if p1.Status != task.Aborted {
		t.Fatal("p1 should have been aborted")
	}

	if p2.Status != task.Failed {
		t.Fatal("p2 should have been failed")
	}

	if p3.Status != task.Aborted {
		t.Fatal("p3 should have been aborted")
	}

	if t2.Status != task.Failed {
		t.Fatal("t2 should have failed")
	}

	if t3.Status != task.Skipped {
		t.Fatal("t3 should have skipped")
	}
}

func TestParallelOutput(t *testing.T) {
	c1 := &task.Task{Name: "c1", Command: "echo a"}
	c2 := &task.Task{Name: "c2", Command: "echo b 1>&2"}
	c3 := &task.Task{Name: "c3", Command: "echo c"}
	c4 := &task.Task{Name: "c1", Command: "echo d 1>&2"}
	c5 := &task.Task{Name: "c2", Command: "echo e"}
	c6 := &task.Task{Name: "c3", Command: "echo f 1>&2"}

	parent := &task.Task{Name: "parent", Parallel: Tasks{c1, c2, c3, c4, c5, c6}}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pipeline{}
	p.runTasks(ctx, cancel, Tasks{parent}, nil)

	for i, v := range []string{"a", "c", "e"} {
		str := strings.Split(parent.Stdout.String(), "\n")[i]
		if str != v {
			t.Fatalf("parent.Stdout should contain %s, not %s", v, str)
		}
	}

	for i, v := range []string{"b", "d", "f"} {
		str := strings.Split(parent.Stderr.String(), "\n")[i]
		if str != v {
			t.Fatalf("parent.Stderr should contain %s, not %s", v, str)
		}
	}

	for i, v := range []string{"a", "b", "c", "d", "e", "f"} {
		str := strings.Split(parent.CombinedOutput.String(), "\n")[i]
		if str != v {
			t.Fatalf("parent.CombinedOutput should contain %s, not %s", v, str)
		}
	}
}

func TestSerialOutput(t *testing.T) {
	c1 := &task.Task{Name: "c1", Command: "echo a"}
	c2 := &task.Task{Name: "c2", Command: "echo b 1>&2"}
	c3 := &task.Task{Name: "c3", Command: "echo c"}
	c4 := &task.Task{Name: "c1", Command: "echo d 1>&2"}
	c5 := &task.Task{Name: "c2", Command: "echo e"}
	c6 := &task.Task{Name: "c3", Command: "echo f && echo g 1>&2"}

	parent := &task.Task{Name: "parent", Serial: Tasks{c1, c2, c3, c4, c5, c6}}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pipeline{}
	p.runTasks(ctx, cancel, Tasks{parent}, nil)

	if !strings.Contains(parent.Stdout.String(), "f") {
		t.Fatal("stdout should contain f")
	}

	if !strings.Contains(parent.Stderr.String(), "g") {
		t.Fatal("stderr should contain g")
	}

	if !strings.Contains(parent.CombinedOutput.String(), "f") {
		t.Fatal("combined output should contain f")
	}

	if !strings.Contains(parent.CombinedOutput.String(), "g") {
		t.Fatal("combined output should contain g")
	}
}

func TestPipe(t *testing.T) {
	t1 := &task.Task{Name: "t1", Command: "echo \"a\nb\""}
	t2 := &task.Task{Name: "t2", Command: "cat"}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pipeline{}
	p.runTasks(ctx, cancel, Tasks{t1, t2}, nil)

	if !strings.Contains(t2.Stdout.String(), "a") {
		t.Fatal("t2.Stdout should contain a")
	}

	if !strings.Contains(t2.Stdout.String(), "b") {
		t.Fatal("t2.Stdout should contain b")
	}
}

func TestPipeOfParallelTasks(t *testing.T) {
	t1 := &task.Task{Name: "t1", Command: "echo t1"}

	p1 := &task.Task{Name: "p1", Command: "cat"}
	p2 := &task.Task{Name: "p2", Command: "cat"}
	p3 := &task.Task{Name: "p3", Command: "echo p3"}

	t2 := &task.Task{Name: "t2", Parallel: Tasks{p1, p2, p3}}

	t3 := &task.Task{Name: "t3", Command: "cat"}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pipeline{}
	p.runTasks(ctx, cancel, Tasks{t1, t2, t3}, nil)

	if !strings.Contains(p1.Stdout.String(), "t1") {
		t.Fatal("p1.Stdout should contain t1")
	}

	if !strings.Contains(p2.Stdout.String(), "t1") {
		t.Fatal("p2.Stdout should contain t1")
	}

	if !strings.Contains(t2.Stdout.String(), "t1") {
		t.Fatal("t2.Stdout should contain t1")
	}

	if !strings.Contains(t2.Stdout.String(), "p3") {
		t.Fatal("t2.Stdout should contain p3")
	}

	if !strings.Contains(t3.Stdout.String(), "t1") {
		t.Fatal("t3.Stdout should contain t1")
	}

	if !strings.Contains(t3.Stdout.String(), "p3") {
		t.Fatal("t3.Stdout should contain p3")
	}
}

func TestPipeOfSerialTasks(t *testing.T) {
	t1 := &task.Task{Name: "t1", Command: "echo t1"}

	s1 := &task.Task{Name: "s1", Command: "cat"}
	s2 := &task.Task{Name: "s2", Command: "cat"}
	s3 := &task.Task{Name: "s3", Command: "echo s3"}

	t2 := &task.Task{Name: "t2", Serial: Tasks{s1, s2, s3}}

	t3 := &task.Task{Name: "t3", Command: "cat"}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pipeline{}
	p.runTasks(ctx, cancel, Tasks{t1, t2, t3}, nil)

	if !strings.Contains(s1.Stdout.String(), "t1") {
		t.Fatal("p1.Stdout should contain t1")
	}

	if !strings.Contains(s2.Stdout.String(), "t1") {
		t.Fatal("p2.Stdout should contain t1")
	}

	if !strings.Contains(t2.Stdout.String(), "s3") {
		t.Fatal("t2.Stdout should contain s3")
	}

	if !strings.Contains(t3.Stdout.String(), "s3") {
		t.Fatal("t3.Stdout should contain s3")
	}
}

func TestExitStatusSuccess(t *testing.T) {
	p := &Pipeline{}
	t1 := &task.Task{Command: "echo"}
	p.Build.Tasks = Tasks{t1}
	code := p.Run(true, true)
	if code != 0 {
		t.Fatalf("Exit code should be 0, not %d", code)
	}
}

func TestExitStatusFail(t *testing.T) {
	p := &Pipeline{}
	t1 := &task.Task{Command: "no_such_command"}
	p.Build.Tasks = Tasks{t1}
	code := p.Run(true, true)
	if code != 1 {
		t.Fatalf("Exit code should be 1, not %d", code)
	}
}

func TestExitStatusParallel(t *testing.T) {
	t1 := &task.Task{Command: "echo"}
	t2 := &task.Task{Command: "no_such_command"}

	p := &Pipeline{}
	p.Build.Tasks = Tasks{&task.Task{Parallel: Tasks{t1, t2}}}
	code := p.Run(true, true)
	if code != 1 {
		t.Fatalf("Exit code should be 1, not %d", code)
	}
}

func TestExitStatusSerial(t *testing.T) {
	t1 := &task.Task{Command: "echo"}
	t2 := &task.Task{Command: "no_such_command"}

	p := &Pipeline{}
	p.Build.Tasks = Tasks{&task.Task{Serial: Tasks{t1, t2}}}
	code := p.Run(true, true)
	if code != 1 {
		t.Fatalf("Exit code should be 1, not %d", code)
	}
}

func TestIncludeInParallel(t *testing.T) {
	tsk := &task.Task{
		Name:     "test include files in parallel task",
		Parallel: []*task.Task{&task.Task{Include: "foo.yml"}},
	}

	p := &Pipeline{}
	p.Build.Tasks = Tasks{tsk}

	code := p.Run(true, true)
	if code != 1 {
		t.Fatalf("Exit code should be 1, not %d", code)
	}
}
