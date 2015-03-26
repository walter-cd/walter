/* walter: a deployment pipeline template
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

	"github.com/recruit-tech/walter/messengers"
	"github.com/recruit-tech/walter/services"
	"github.com/recruit-tech/walter/stages"
)

type Pipeline struct {
	Stages list.List
}

type Resources struct {
	Pipeline    *Pipeline
	Cleanup     *Pipeline
	Reporter    messengers.Messenger
	RepoService services.Service
}

func (self *Resources) ReportStageResult(stage stages.Stage, result bool) {
	name := stage.GetStageName()
	self.Reporter.Post(
		fmt.Sprintf("Stage execution results: %+v, %+v", name, result))

	if stage.GetStageOpts().ReportingFullOutput {
		if out := stage.GetOutResult(); len(out) > 0 {
			self.Reporter.Post(
				fmt.Sprintf("[%s] %s", name, stage.GetOutResult()))
		}
		if err := stage.GetErrResult(); len(err) > 0 {
			self.Reporter.Post(
				fmt.Sprintf("[%s][ERROR] %s", name, stage.GetErrResult()))
		}
	}
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
