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

// Package pipelines defines the pipeline and the resources.
package pipelines

import (
	"container/list"
	"fmt"

	"github.com/recruit-tech/walter/messengers"
	"github.com/recruit-tech/walter/services"
	"github.com/recruit-tech/walter/stages"
)

// Pipeline stores the list of stages.
type Pipeline struct {
	Stages list.List
}

// Resources stores the settings loaded from the configuation file.
type Resources struct {
	// Pipeline stores the list of stages to be executed.
	Pipeline *Pipeline

	// Cleanup stores the list of stages extecuted after Pipeline.
	Cleanup *Pipeline

	// Reporter stores the messenger client which report the result to the server.
	Reporter messengers.Messenger

	// RepoService is a client of VCS services such as GitHub and reports the result to the service.
	RepoService services.Service
}

// ReportStageResult throw the results of specified stage to the messenger services.
func (self *Resources) ReportStageResult(stage stages.Stage, result string) {
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

// AddStage appends specified stage to the pipeline.
func (self *Pipeline) AddStage(stage stages.Stage) {
	self.Stages.PushBack(stage)
}

// Size returns the number of stages in the pipeline.
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
