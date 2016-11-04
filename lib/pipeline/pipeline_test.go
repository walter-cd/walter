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
	p.runTasks(ctx, cancel, tasks)

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
	p.runTasks(ctx, cancel, tasks)

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
	p.runTasks(ctx, cancel, Tasks{parent})

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
