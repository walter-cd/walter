package plumber

import "container/list"

type Pipeline struct {
	stages list.List
}

func (self *Pipeline) Run() {
	// TODO: apply dependency 
	for stage := self.stages.Front(); stage != nil; stage = stage.Next() {
		// do something with e.Value
	}
}

func (self *Pipeline) AddStage(stage Stage) {
	self.stages.PushBack(stage)
}

func (self *Pipeline) Size() int {
	return self.stages.Len()
}

func  NewPipeline() *Pipeline {
	return &Pipeline{}
}
