package pipeline

import (
	"container/list"
	"fmt"

	"github.com/takahi-i/plumber/stage"
)

type Pipeline struct {
	stages list.List
}

func (self *Pipeline) Run() bool {
	// TODO: apply dependency
	for stageItem := self.stages.Front(); stageItem != nil; stageItem = stageItem.Next() {
		fmt.Printf("Executing planned stage: %s\n", stageItem.Value)
		stageItem.Value.(stage.Stage).Run()
	}
	return true
}

func (self *Pipeline) AddStage(stage stage.Stage) {
	self.stages.PushBack(stage)
}

func (self *Pipeline) Size() int {
	return self.stages.Len()
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}
