package pipeline

import (
	"fmt"
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
