package plumber

import (
	"testing"
)

func TestSum(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddStage(NewCommandStage())
	pipeline.AddStage(NewCommandStage())
	expected := 2
	actual := pipeline.Size()
	if 2 != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
