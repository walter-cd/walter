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

// Package stages defines the stages which are included in pipelines.
package stages

import (
	"container/list"

	"github.com/walter-cd/walter/log"
)

// BaseStage is an abstract struct implemented by the inherited struct that wishes to run somthing
// in a stage.
type BaseStage struct {
	// Run specified process. The required method (Run) is defined in the inheirted structs.
	Runner

	// Input channel. This channel stores the input of the statge.
	InputCh *chan Mediator

	// Output channel. This channel stores the output of the statge.
	OutputCh *chan Mediator

	// List of child stages.
	ChildStages list.List

	// Stage name.
	StageName string `config:"name"`

	// Results of stdout flush by the stage.
	OutResult string

	// Results of stderr flush by the stage.
	ErrResult string

	// Return value of the stage
	ReturnValue bool

	// options of the stage.
	Opts StageOpts
}

//StageOpts struct for handing stage outputs
type StageOpts struct {
	// Flush all output when the value is true.
	ReportingFullOutput bool `config:"report_full_output"`
}

//NewStageOpts creates a new stage output
func NewStageOpts() *StageOpts {
	return &StageOpts{
		ReportingFullOutput: false,
	}
}

// Run executes the stage.
func (b *BaseStage) Run() bool {
	if b.Runner == nil {
		panic("Mast have a child class assigned")
	}
	b.ReturnValue = b.Runner.Run()
	return b.ReturnValue
}

// AddChildStage appends one child stage.
func (b *BaseStage) AddChildStage(stage Stage) {
	log.Debugf("added childstage: %v", stage)
	b.ChildStages.PushBack(stage)
}

// GetChildStages returns a list of child stages.
func (b *BaseStage) GetChildStages() list.List {
	return b.ChildStages
}

//GetStageName returns the name of the current stage
func (b *BaseStage) GetStageName() string {
	return b.StageName
}

// SetStageName sets stage name.
func (b *BaseStage) SetStageName(stageName string) {
	b.StageName = stageName
}

// GetStageOpts returns stage options.
func (b *BaseStage) GetStageOpts() StageOpts {
	return b.Opts
}

// SetStageOpts sets stage options.
func (b *BaseStage) SetStageOpts(stageOpts StageOpts) {
	b.Opts = stageOpts
}

// SetInputCh sets input channel.
func (b *BaseStage) SetInputCh(inputCh *chan Mediator) {
	b.InputCh = inputCh
}

// GetInputCh retruns input channel.
func (b *BaseStage) GetInputCh() *chan Mediator {
	return b.InputCh
}

// SetOutputCh sets output channel.
func (b *BaseStage) SetOutputCh(outputCh *chan Mediator) {
	b.OutputCh = outputCh
}

// GetOutputCh retruns output channel.
func (b *BaseStage) GetOutputCh() *chan Mediator {
	return b.OutputCh
}

// GetOutResult returns standard output results.
func (b *BaseStage) GetOutResult() string {
	return b.OutResult
}

// SetOutResult sets standard output results.
func (b *BaseStage) SetOutResult(result string) {
	b.OutResult = result
}

// GetErrResult returns standard error results.
func (b *BaseStage) GetErrResult() string {
	return b.ErrResult
}

// SetErrResult sets standard error results.
func (b *BaseStage) SetErrResult(result string) {
	b.ErrResult = result
}

// GetReturnValue returns return value of the stage
func (b *BaseStage) GetReturnValue() bool {
	return b.ReturnValue
}

// SetReturnValue sets return value of the stage
func (b *BaseStage) SetReturnValue(value bool) {
	b.ReturnValue = value
}
