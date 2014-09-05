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

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
)

type CommandStage struct {
	BaseStage
	Command   string   `config:"command"`
	Arguments []string `config:"arguments"`
	OutResult string
}

func (self *CommandStage) GetStdoutResult() string {
	return self.OutResult
}

func (self *CommandStage) Run() bool {
	command := strings.Split(self.Command, " ")
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Args = append(command, self.Arguments...)
	cmd.Dir = "."
	out, err := cmd.StdoutPipe()

	if err != nil {
		return false
	}
	err = cmd.Start()
	if err != nil {
		return false
	}
	self.OutResult = copyStream(out)
	err = cmd.Wait()
	if err != nil {
		return false
	}
	return true
}

func copyStream(reader io.Reader) string {
	var err error
	var n int
	var buffer bytes.Buffer
	tmpBuf := make([]byte, 1024)
	for {
		if n, err = reader.Read(tmpBuf); err != nil {
			break
		}
		buffer.Write(tmpBuf[0:n])
	}
	if err == io.EOF {
		err = nil
	} else {
		log.Println("ERROR: " + err.Error())
	}
	return buffer.String()
}

func (self *CommandStage) AddCommand(command string, arguments ...string) {
	self.Command = command
	self.Arguments = arguments
	self.BaseStage = BaseStage{}
	self.BaseStage.Runner = self
}

func NewCommandStage() *CommandStage {
	stage := CommandStage{}
	return &stage
}
