package stage

import (
	"testing"
)

func TestWIthSimpleCommand(t *testing.T) {
	stage := NewCommandStage()
	stage.AddCommand("ls", "-l")
	expected := true
	actual := stage.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestWithNoexistCommand(t *testing.T) {
	stage := NewCommandStage()
	stage.AddCommand("zzzz", "")
	expected := false
	actual := stage.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
