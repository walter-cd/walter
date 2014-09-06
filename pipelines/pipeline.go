/* plumber: a deployment pipeline template
 * Copyright (C) 2014 Recruit Technologies Co., Ltd. and contributors
 * (see CONTRIBUTORS.md)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package pipelines

import (
	"container/list"

	"github.com/recruit-tech/plumber/log"
	"github.com/recruit-tech/plumber/stages"
)

type Pipeline struct {
	Stages list.List
}

func (self *Pipeline) Run() bool {
	mediator := stages.Mediator{States: make(map[string]string)}
	log.Info("geting starting to run pipeline process...")
	for stageItem := self.Stages.Front(); stageItem != nil; stageItem = stageItem.Next() {
		log.Debugf("Executing planned stage: %s\n", stageItem.Value)
		mediator = stages.Execute(stageItem.Value.(stages.Stage), mediator)
	}
	log.Info("finished to run pipeline process...")
	return true
}

func (self *Pipeline) AddStage(stage stages.Stage) {
	self.Stages.PushBack(stage)
}

func (self *Pipeline) Size() int {
	return self.Stages.Len()
}

func (self *Pipeline) Build() {
	self.buildDeps(&self.Stages)
}

func (self *Pipeline) buildDeps(stages *list.List) {
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}
