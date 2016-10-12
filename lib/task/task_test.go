package task

import (
	"strings"
	"testing"

	"golang.org/x/net/context"
)

func TestStdout(t *testing.T) {
	tsk := &Task{Name: "echo", Command: "echo hello"}

	ctx, cancel := context.WithCancel(context.Background())

	err := tsk.Run(ctx, cancel)

	if err != nil {
		t.Fatal(err)
	}

	if !contains(tsk.Stdout, "hello") {
		t.Fatalf("tsk.Stdout is %s, it does not contain \"hello\"", tsk.Stdout)
	}

	if !contains(tsk.CombinedOutput, "hello") {
		t.Fatalf("tsk.CombinedOutput is %s, it does not contain \"hello\"", tsk.CombinedOutput)
	}
}

func TestStderr(t *testing.T) {
	tsk := &Task{Name: "echo", Command: "echo hello 1>&2"}

	ctx, cancel := context.WithCancel(context.Background())

	err := tsk.Run(ctx, cancel)

	if err != nil {
		t.Fatal(err)
	}

	if !contains(tsk.Stderr, "hello") {
		t.Fatalf("tsk.Stderr is %s, it does not contain \"hello\"", tsk.Stderr)
	}

	if !contains(tsk.CombinedOutput, "hello") {
		t.Fatalf("tsk.CombinedOutput is %s, it does not contain \"hello\"", tsk.CombinedOutput)
	}
}

func TestStatus(t *testing.T) {
	tsk := &Task{Name: "command should succeed", Command: "echo foo"}

	ctx, cancel := context.WithCancel(context.Background())
	err := tsk.Run(ctx, cancel)
	if err != nil {
		t.Fatal(err)
	}
	if tsk.Status != Succeeded {
		t.Fatal("command not succeeded")
	}

	tsk = &Task{Name: "command should fail", Command: "no_such_command"}
	err = tsk.Run(ctx, cancel)
	if err != nil {
		t.Fatal(err)
	}

	if tsk.Status != Failed {
		t.Fatal("command not failed")
	}
}

func TestSerialTasks(t *testing.T) {
	t1 := &Task{Name: "foo", Command: "echo foo"}
	t2 := &Task{Name: "bar", Command: "barbarbar"}
	t3 := &Task{Name: "baz", Command: "echo baz"}

	tasks := &Tasks{t1, t2, t3}

	ctx, cancel := context.WithCancel(context.Background())
	tasks.Run(ctx, cancel)

	if t1.Status != Succeeded {
		t.Fatal("t1 should have succeeded")
	}

	if t2.Status != Failed {
		t.Fatal("t2 should have failed")
	}

	if t3.Status != Skipped {
		t.Fatalf("t2 should have beed skipped")
	}
}

func TestParallelTasks(t *testing.T) {
	p1 := &Task{Name: "p1", Command: "sleep 1"}
	p2 := &Task{Name: "p2", Command: "p2p2p2p2"}
	p3 := &Task{Name: "p3", Command: "sleep 1"}

	t1 := &Task{Name: "foo", Command: "echo foo"}
	t2 := &Task{Name: "bar", Parallel: Parallel{p1, p2, p3}}
	t3 := &Task{Name: "baz", Command: "echo baz"}

	tasks := &Tasks{t1, t2, t3}
	ctx, cancel := context.WithCancel(context.Background())
	tasks.Run(ctx, cancel)

	if p1.Status != Aborted {
		t.Fatal("p1 should have been aborted")
	}

	if p2.Status != Failed {
		t.Fatal("p2 should have been failed")
	}

	if p3.Status != Aborted {
		t.Fatal("p3 should have been aborted")
	}

	if t2.Status != Failed {
		t.Fatal("t2 should have failed")
	}

	if t3.Status != Skipped {
		t.Fatal("t3.should have skipped")
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Contains(a, e) {
			return true
		}
	}
	return false
}
