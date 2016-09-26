package task

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tsk := Task{Name: "echo", Command: "echo hello"}

	err := tsk.Run()

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(tsk.Stdout, "hello") {
		t.Fatalf("tsk.Stdout is %s, it does not contain \"hello\"", tsk.Stdout)
	}
}
