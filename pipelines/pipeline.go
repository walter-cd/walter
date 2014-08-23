package pipelines

import (
	"container/list"
	"fmt"

	"github.com/recruit-tech/plumber/stages"
)

type Pipeline struct {
	Stages list.List
}

func (self *Pipeline) Run() bool {
	// TODO: apply dependency
	for stageItem := self.Stages.Front(); stageItem != nil; stageItem = stageItem.Next() {
		fmt.Printf("Executing planned stage: %s\n", stageItem.Value)
		stageItem.Value.(stages.Stage).Run()
	}
	return true
}

func (self *Pipeline) AddStage(stage stages.Stage) {
	self.Stages.PushBack(stage)
}

func (self *Pipeline) Size() int {
	return self.Stages.Len()
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}
