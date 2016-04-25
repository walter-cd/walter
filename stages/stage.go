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

//Package stages contains functionality for managing stage lifecycle
package stages

import (
	"container/list"
	"fmt"
)

// Stage is a interface type which declares a list of methods every Stage object should define.
type Stage interface {
	AddChildStage(Stage)
	GetChildStages() list.List
	GetStageName() string
	SetStageName(string)
	GetStageOpts() StageOpts
	SetStageOpts(StageOpts)
	GetInputCh() *chan Mediator
	SetInputCh(*chan Mediator)
	GetOutputCh() *chan Mediator
	SetOutputCh(*chan Mediator)
	GetOutResult() string
	SetOutResult(string)
	GetErrResult() string
	SetErrResult(string)
	GetCombinedResult() string
	SetCombinedResult(string)
	GetSuppressAll() bool
	SetSuppressAll(bool)
	GetReturnValue() bool
}

// Runner contains the Run method which is deined in Stage implemantations.
type Runner interface {
	Run() bool
}

// Mediator stores the intermidate results.
type Mediator struct {
	States map[string]string
	Type   string
}

// IsAnyFailure returns true when mediator found any failures. Otherwise This method returns false.
func (m *Mediator) IsAnyFailure() bool {
	for _, v := range m.States {
		if v == "false" {
			return true
		}
	}
	return false
}

// InitStage initializes a stage with specified stage type.
func InitStage(stageType string) (Stage, error) {
	var stage Stage
	switch stageType {
	case "command":
		stage = new(CommandStage)
	case "shell":
		stage = new(ShellScriptStage)
	default:
		return nil, fmt.Errorf("No specified stage type: '%s'", stageType)
	}
	PrepareCh(stage)
	return stage, nil
}

// PrepareCh prepares input and output channels.
func PrepareCh(stage Stage) {
	in := make(chan Mediator)
	out := make(chan Mediator)
	stage.SetInputCh(&in)
	stage.SetOutputCh(&out)
}
