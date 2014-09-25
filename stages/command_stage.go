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

type CommandStage struct {
	BaseStage
	Command   string `config:"command"`
	OutResult string
	ErrResult string
}

func (self *CommandStage) GetStdoutResult() string {
	return self.OutResult
}

func (self *CommandStage) Run() bool {
	cmd := exec.Command("sh", "-c", self.Command)
	log.Infof("[command] exec: %s", self.Command)
	cmd.Dir = "."
	out, err := cmd.StdoutPipe()
	outE, errE := cmd.StderrPipe()

	if err != nil {
		log.Errorf("[command] err: %s", out)
		log.Errorf("[command] err: %s", err)
		return false
	}

	if errE != nil {
		log.Errorf("[command] err: %s", outE)
		log.Errorf("[command] err: %s", errE)
		return false
	}

	err = cmd.Start()
	if err != nil {
		log.Errorf("[command] err: %s", err)
		return false
	}
	self.OutResult = copyStream(out)
	self.ErrResult = copyStream(outE)

	err = cmd.Wait()
	if err != nil {
		log.Errorf("[command] err: %s", err)
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
		log.Infof("[command] output: %s", tmpBuf[0:n])
	}
	if err == io.EOF {
		err = nil
	} else {
		log.Error("ERROR: " + err.Error())
	}
	return buffer.String()
}

func (self *CommandStage) AddCommand(command string) {
	self.Command = command
	self.BaseStage.Runner = self
}

func NewCommandStage() *CommandStage {
	stage := CommandStage{}
	return &stage
}
