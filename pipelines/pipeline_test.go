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
	"testing"

	"github.com/walter-cd/walter/stages"
	"github.com/stretchr/testify/assert"
)

func createStage(stageType string) stages.Stage {
	stage, _ := stages.InitStage(stageType)
	return stage
}

func createCommandStage(command string) *stages.CommandStage {
	in := make(chan stages.Mediator)
	out := make(chan stages.Mediator)
	return &stages.CommandStage{
		Command: "echo",
		BaseStage: stages.BaseStage{
			InputCh:  &in,
			OutputCh: &out,
		},
	}
}

func TestAddPipeline(t *testing.T) {
	pipeline := NewPipeline()
	pipeline.AddStage(createStage("command"))
	pipeline.AddStage(createStage("command"))
	assert.Equal(t, 2, pipeline.Size())
}

type MockMessenger struct {
	Posts []string
}

func (mock *MockMessenger) Post(msg string) bool {
	mock.Posts = append(mock.Posts, msg)
	return true
}

func (mock *MockMessenger) Suppress(outputType string) bool {
	return false
}

func TestReportStageResult(t *testing.T) {
	mock := &MockMessenger{}
	p := Resources{
		Reporter: mock,
	}

	stage := createStage("command")
	stage.SetStageName("test")

	opts := stages.NewStageOpts()

	stage.SetStageOpts(*opts)

	p.ReportStageResult(stage, "true")

	assert.Equal(t, 1, len(mock.Posts))
}

func TestReportStageResultWithFullOutput(t *testing.T) {
	mock := &MockMessenger{}
	p := Resources{
		Reporter: mock,
	}

	stage := createStage("command")
	stage.SetStageName("test")
	stage.SetOutResult("output")

	opts := stages.NewStageOpts()
	opts.ReportingFullOutput = true

	stage.SetStageOpts(*opts)

	p.ReportStageResult(stage, "true")

	assert.Equal(t, 2, len(mock.Posts))
}
