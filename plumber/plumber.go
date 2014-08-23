package plumber

import (
	"github.com/recruit-tech/plumber/config"
	"github.com/recruit-tech/plumber/pipelines"
)

type Plumber struct {
	Pipeline *pipelines.Pipeline
}

func New(opts *config.Opts) *Plumber {
	configData := config.ReadConfig(opts.PipelineFilePath)
	pipeline := (config.Parse(configData))
	return &Plumber{
		Pipeline: pipeline,
	}
}

func (e *Plumber) Run() {
	e.Pipeline.Run()
}
