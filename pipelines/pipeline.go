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
func (resources *Resources) ReportStageResult(stage stages.Stage, resultStr string) {
	name := stage.GetStageName()
	if !resources.Reporter.Suppress("result") {
		if resultStr == "true" {
			resources.Reporter.Post(
				fmt.Sprintf("[%s][RESULT] Succeeded", name))
		} else if resultStr == "skipped" {
			resources.Reporter.Post(
				fmt.Sprintf("[%s][RESULT] Skipped", name))
		} else {
			resources.Reporter.Post(
				fmt.Sprintf("[%s][RESULT] Failed", name))
		}
	}

	if stage.GetStageOpts().ReportingFullOutput {
		if out := stage.GetOutResult(); (len(out) > 0) && (!resources.Reporter.Suppress("stdout")) {
			resources.Reporter.Post(
				fmt.Sprintf("[%s][STDOUT] %s", name, stage.GetOutResult()))
		}
		if err := stage.GetErrResult(); len(err) > 0 && (!resources.Reporter.Suppress("stderr")) {
			resources.Reporter.Post(
				fmt.Sprintf("[%s][STDERR] %s", name, stage.GetErrResult()))
		}
	}
}

// AddStage appends specified stage to the pipeline.
func (resources *Pipeline) AddStage(stage stages.Stage) {
	resources.Stages.PushBack(stage)
}

// GetStageResult returns the result (stdout, stderr, return value) of specified stage.
func (resources *Pipeline) GetStageResult(name string, stageType string) (string, error) {
	for stageItem := resources.Stages.Front(); stageItem != nil; stageItem = stageItem.Next() {
		stage := stageItem.Value.(stages.Stage)
		if name != stage.GetStageName() {
			continue
		}
		switch stageType {
		case "__OUT":
			return stage.GetErrResult(), nil
		case "__ERR":
			return stage.GetOutResult(), nil
		case "__RESULT":
			return "0", nil // TODO: fixme
		default:
			return "", fmt.Errorf("no specified type: " + stageType)
		}
	}
	return "", fmt.Errorf("no specified stage: " + name)
}

// Size returns the number of stages in the pipeline.
func (resources *Pipeline) Size() int {
	return resources.Stages.Len()
}

//Build builds a pipeline for the current resources
func (resources *Pipeline) Build() {
	resources.buildDeps(&resources.Stages)
}

func (resources *Pipeline) buildDeps(stages *list.List) {
}

//NewPipeline create a new pipeline instance
func NewPipeline() *Pipeline {
	return &Pipeline{}
}
