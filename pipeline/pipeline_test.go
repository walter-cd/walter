package pipeline

import (
	"testing"

	"github.com/takahi-i/plumber/stage"
)

func createCommandStage(command string, arguments ...string) *stage.CommandStage {
	resultStage := stage.NewCommandStage()
	resultStage.AddCommand(command, arguments...)
	return resultStage
}

func TestAddPipeline(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddStage(stage.NewCommandStage())
	pipeline.AddStage(stage.NewCommandStage())
	expected := 2
	actual := pipeline.Size()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestRunPipeline(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddStage(createCommandStage("echo", "foobar"))
	pipeline.AddStage(createCommandStage("echo", "baz"))
	expected := true
	actual := pipeline.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
	// TODO: check the output from stage
}
