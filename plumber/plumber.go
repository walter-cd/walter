package plumber

import (
	"github.com/takahi-i/plumber/pipeline"
	"github.com/takahi-i/plumber/stage"
)

type Plumber struct {
	Pipeline *pipeline.Pipeline
}

func New() *Plumber {
	var pipeline = pipeline.NewPipeline()
	return &Plumber{
		Pipeline: pipeline,
	}
}

func (e *Plumber) Run() {
	var stage = stage.NewCommandStage()
	stage.AddCommand("echo", "Hello, I'm command stage'")

	e.Pipeline.AddStage(stage)
	e.Pipeline.Run()
}
