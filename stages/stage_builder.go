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

//Builder object to construct a Stage object.
type StageBuilder struct {
	stage Stage
}

//Create and initalize stage.
func (b *StageBuilder) NewStage(stageType string) *StageBuilder {
	b.stage, _ = InitStage(stageType)
	return b
}

//Set stdout result (this method is for testing).
func (b *StageBuilder) SetOutResult(out string) *StageBuilder {
	b.stage.SetOutResult(out)
	return b
}

//Set stdout result (this method is for testing).
func (b *StageBuilder) SetErrResult(err string) *StageBuilder {
	b.stage.SetOutResult(err)
	return b
}

//Set stage name(this method is for testing).
func (b *StageBuilder) SetName(name string) *StageBuilder {
	b.stage.SetStageName(name)
	return b
}

//Set stage target (this method is for testing).
func (b *StageBuilder) SetTarget(target string) *StageBuilder {
	switch b.stage.(type) {
	case *CommandStage:
		command := b.stage.(*CommandStage)
		command.AddCommand(target)
	case *ShellScriptStage:
		shell := b.stage.(*ShellScriptStage)
		shell.File = target
	}
	return b
}

//Return the initalized stage.
func (b *StageBuilder) Build() Stage {
	return b.stage
}
