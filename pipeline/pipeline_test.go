package pipeline

import (
	"testing"

	"github.com/takahi-i/plumber/stage"
)

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
