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
