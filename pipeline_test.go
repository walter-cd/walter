package plumber

import (
	"testing"
)

func TestAddPipeline(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddStage(NewCommandStage())
	pipeline.AddStage(NewCommandStage())
	expected := 2
	actual := pipeline.Size()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
