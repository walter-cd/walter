package task

import (
	"strings"
	"testing"
)

func TestStdout(t *testing.T) {
	tsk := Task{Name: "echo", Command: "echo hello"}

	err := tsk.Run()

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
	tsk := Task{Name: "echo", Command: "echo hello 1>&2"}

	err := tsk.Run()

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
	tsk := Task{Name: "command should succeed", Command: "echo foo"}
	err := tsk.Run()
	if err != nil {
		t.Fatal(err)
	}
	if tsk.Status != Succeeded {
		t.Fatal("command not succeeded")
	}

	tsk = Task{Name: "command should fail", Command: "no_such_command"}
	err = tsk.Run()
	if err != nil {
		t.Fatal(err)
	}
	if tsk.Status != Failed {
		t.Fatal("command not failed")
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
