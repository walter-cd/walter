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

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/recruit-tech/walter/log"
)

// CommandStage executes more than one commands.
type CommandStage struct {
	BaseStage
	Command   string `config:"command" is_replace:"false"`
	Directory string `config:"directory" is_replace:"true"`
	OnlyIf    string `config:"only_if" is_replace:"false"`
}

// Get standard output results.
func (self *CommandStage) GetStdoutResult() string {
	return self.OutResult
}

// Run registered commands.
func (self *CommandStage) Run() bool {
	// Check OnlyIf
	if self.runOnlyIf() == false {
		log.Warnf("[command] exec: skipped this stage \"%s\", since only_if condition failed", self.BaseStage.StageName)
		return true
	}

	// Run command
	result := self.runCommand()
	if result == false {
		log.Errorf("[command] exec: failed stage \"%s\"", self.BaseStage.StageName)
	}
	return result
}

func (self *CommandStage) runOnlyIf() bool {
	if self.OnlyIf == "" {
		return true
	}
	log.Infof("[command] only_if: found \"only_if\" attribute in stage \"%s\"", self.BaseStage.StageName)
	cmd := exec.Command("sh", "-c", self.OnlyIf)
	log.Infof("[command] only_if: %s", self.BaseStage.StageName)
	log.Debugf("[command] only_if literal: %s", self.OnlyIf)
	cmd.Dir = self.Directory
	result, _, _ := execCommand(cmd, "only_if", self.BaseStage.StageName)
	return result
}

func (self *CommandStage) runCommand() bool {
	cmd := exec.Command("sh", "-c", self.Command)
	log.Infof("[command] exec: %s", self.BaseStage.StageName)
	log.Debugf("[command] exec command literal: %s", self.Command)
	cmd.Dir = self.Directory
	result, outResult, errResult := execCommand(cmd, "exec", self.BaseStage.StageName)
	self.SetOutResult(*outResult)
	self.SetErrResult(*errResult)
	return result
}

func execCommand(cmd *exec.Cmd, prefix string, name string) (bool, *string, *string) {
	out, err := cmd.StdoutPipe()
	outE, errE := cmd.StderrPipe()

	if err != nil {
		log.Warnf("[command] %s err: %s", prefix, out)
		log.Warnf("[command] %s err: %s", prefix, err)
		return false, nil, nil
	}

	if errE != nil {
		log.Warnf("[command] %s err: %s", prefix, outE)
		log.Warnf("[command] %s err: %s", prefix, errE)
		return false, nil, nil
	}

	err = cmd.Start()
	if err != nil {
		log.Warnf("[command] %s err: %s", prefix, err)
		return false, nil, nil
	}
	outResult := copyStream(out, prefix, name)
	errResult := copyStream(outE, prefix, name)

	err = cmd.Wait()
	if err != nil {
		log.Warnf("[command] %s err: %s", prefix, err)
		return false, &outResult, &errResult
	}
	return true, &outResult, &errResult
}

func copyStream(reader io.Reader, prefix string, name string) string {
	var err error
	var n int
	var buffer bytes.Buffer
	tmpBuf := make([]byte, 1024)
	for {
		if n, err = reader.Read(tmpBuf); err != nil {
			break
		}
		buffer.Write(tmpBuf[0:n])
		log.Infof("[%s][command] %s output: %s", name, prefix, tmpBuf[0:n])
	}
	if err == io.EOF {
		err = nil
	} else {
		log.Error("ERROR: " + err.Error())
	}
	return buffer.String()
}

// Register specified command.
func (self *CommandStage) AddCommand(command string) {
	self.Command = command
	self.BaseStage.Runner = self
}

// Set the directory where the command is executed.
func (self *CommandStage) SetDirectory(directory string) {
	self.Directory = directory
}

// Create one CommandStage object.
func NewCommandStage() *CommandStage {
	stage := CommandStage{Directory: "."}
	return &stage
}
