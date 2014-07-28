package plumber

import (
	"testing"
)

func Test(t *testing.T) {
	stage := NewCommandStage()
	stage.AddCommand("ls", []string{"-l"})
	expected := true
	actual := stage.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
