package task

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

func TestStdout(t *testing.T) {
	tsk := &Task{Name: "echo", Command: "echo hello"}

	ctx, cancel := context.WithCancel(context.Background())

	err := tsk.Run(ctx, cancel, nil)

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

	err := tsk.Run(ctx, cancel, nil)

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
	err := tsk.Run(ctx, cancel, nil)
	if err != nil {
		t.Fatal(err)
	}
	if tsk.Status != Succeeded {
		t.Fatal("command not succeeded")
	}

	tsk = &Task{Name: "command should fail", Command: "no_such_command"}
	err = tsk.Run(ctx, cancel, nil)
	if err == nil {
		t.Fatal("tsk.Run() should return err")
	}

	if tsk.Status != Failed {
		t.Fatal("tsk.Status shoud be Fail")
	}
}

func contains(buf *bytes.Buffer, e string) bool {
	for _, a := range strings.Split(buf.String(), "\n") {
		if strings.Contains(a, e) {
			return true
		}
	}
	return false
}
