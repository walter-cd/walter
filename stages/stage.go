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
package stages

import "container/list"

type Stage interface {
	AddChildStage(Stage)
	GetChildStages() list.List
	GetStageName() string
	SetStageName(string)
	GetInputCh() *chan Mediator
	SetInputCh(*chan Mediator)
	GetOutputCh() *chan Mediator
	SetOutputCh(*chan Mediator)
}

type Runner interface {
	Run() bool
}

type Mediator struct {
	States map[string]string
	Type   string
}

func InitStage(stageType string) Stage {
	var stage Stage
	switch stageType {
	case "command":
		stage = new(CommandStage)
	case "shell":
		stage = new(ShellScriptStage)
	}
	PrepareCh(stage)
	return stage
}

func PrepareCh(stage Stage) {
	in := make(chan Mediator)
	out := make(chan Mediator)
	stage.SetInputCh(&in)
	stage.SetOutputCh(&out)
}
