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
	"os"
	"os/exec"

	"github.com/walter-cd/walter/log"
)

// ResourceValidator class check if the resources to run the target staget are satisfied.
type ResourceValidator struct {
	files   list.List
	command string
}

//Validate validates if the command can be executed
func (resourceValidator *ResourceValidator) Validate() bool {
	// check if files exists
	for file := resourceValidator.files.Front(); file != nil; file = file.Next() {
		filePath := file.Value.(string)
		log.Debugf("checking file: %v", filePath)
		if _, err := os.Stat(filePath); err == nil {
			log.Debugf("file exists")
		} else {
			log.Errorf("file: %v does not exists", filePath)
			return false
		}
	}
	// check if command exists
	if len(resourceValidator.command) == 0 { // return true when no command is registrated
		return true
	}
	cmd := exec.Command("which", resourceValidator.command)
	err := cmd.Run()
	if err != nil {
		log.Errorf("command: %v does not exists", resourceValidator.command)
		return false
	}
	return true
}

//AddFile add the suplied file to the validator file list
// TODO add permission
func (resourceValidator *ResourceValidator) AddFile(f string) {
	resourceValidator.files.PushBack(f)
}

//AddCommandName adds the command to the validator
func (resourceValidator *ResourceValidator) AddCommandName(c string) {
	resourceValidator.command = c
}

//NewResourceValidator creates a new ResourceValidator
func NewResourceValidator() *ResourceValidator {
	validator := ResourceValidator{}
	return &validator
}
