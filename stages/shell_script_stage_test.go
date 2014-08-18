package stages

import (
	"testing"
)

func TestWithSimpleScript(t *testing.T) {
	stage := NewShellScriptStage()
	stage.AddScript("test_sample.sh")
	expected := true
	actual := stage.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestWithNonExistScript(t *testing.T) {
	stage := NewShellScriptStage()
	stage.AddScript("non_exist_sample.sh")
	expected := false
	actual := stage.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
