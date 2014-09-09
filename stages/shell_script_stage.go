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
package stages

type ShellScriptStage struct {
	command CommandStage
}

func (self *ShellScriptStage) Run() bool {
	return self.command.Run()
}

func (self *ShellScriptStage) AddScript(scriptFile string) {
	// TODO: validate the existance of scriptFile
	// and flush log when the file does not exist.
	self.command.AddCommand("sh " + scriptFile)
}

func NewShellScriptStage() *ShellScriptStage {
	return &ShellScriptStage{}
}
