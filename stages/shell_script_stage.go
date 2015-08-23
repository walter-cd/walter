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
	"github.com/recruit-tech/walter/log"
)

// ShellScriptStage executes one shell script file.
type ShellScriptStage struct {
	ResourceValidator
	CommandStage
	File string `config:"file"`
}

func (shellScriptStage *ShellScriptStage) preCheck() bool {
	shellScriptStage.AddCommandName("sh")
	shellScriptStage.AddFile(shellScriptStage.File)
	return shellScriptStage.Validate()
}

// Run exectues specified shell script.
func (shellScriptStage *ShellScriptStage) Run() bool {
	log.Infof("[shell] exec: %s", shellScriptStage.BaseStage.StageName)
	log.Debugf("[shell] specified file: %s\n", shellScriptStage.File)
	if shellScriptStage.preCheck() == false {
		log.Infof("failed preCheck before running script...")
		return false
	}
	shellScriptStage.AddCommand("sh " + shellScriptStage.File)
	return shellScriptStage.CommandStage.Run()
}

//NewShellScriptStage generate one ShellScriptStage object.
func NewShellScriptStage() *ShellScriptStage {
	return &ShellScriptStage{}
}
