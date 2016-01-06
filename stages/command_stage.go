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
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/walter-cd/walter/log"
)

// CommandStage executes more than one commands.
type CommandStage struct {
	BaseStage
	Command   string `config:"command" is_replace:"false"`
	Directory string `config:"directory" is_replace:"true"`
	OnlyIf    string `config:"only_if" is_replace:"false"`
	WaitFor   string `config:"wait_for" is_replace:"true"`
}

// WaitFor wait until the predefined conditions are satisfied
type WaitFor struct {
	Host  string
	Port  int
	File  string
	State string
	Delay float64
}

//GetStdoutResult returns the stdio output from the command.
func (waitFor *WaitFor) Wait() {
	time.Sleep(100 * time.Millisecond)
}

func ParseWaitFor(waitForStr string) (*WaitFor, error) {
	var wait = &WaitFor{}
	for _, seg := range strings.Split(waitForStr, " ") {
		kv := strings.Split(seg, "=")
		if len(kv) != 2 {
			return nil, errors.New("Given segment does not have two segments: " + seg)
		}
		switch kv[0] {
		case "host":
			wait.Host = kv[1]
		case "state":
			wait.State = kv[1]
		case "port":
			v, err := strconv.Atoi(kv[1])
			if err != nil {
				return nil, err
			}
			wait.Port = v
		case "delay":
			v, err := strconv.ParseFloat(kv[1], 64)
			if err != nil {
				return nil, err
			}
			wait.Delay = v
		case "file":
			wait.File = kv[1]
		default:
			return nil, errors.New("No wait_for property such as: " + kv[0])
		}
	}
	return wait, nil
}

//GetStdoutResult returns the stdio output from the command.
func (commandStage *CommandStage) GetStdoutResult() string {
	return commandStage.OutResult
}

// Run registered commands.
func (commandStage *CommandStage) Run() bool {
	// Check WaitFor
	commandStage.waitFor()

	// Check OnlyIf
	if commandStage.runOnlyIf() == false {
		log.Warnf("[command] exec: skipped this stage \"%s\", since only_if condition failed", commandStage.BaseStage.StageName)
		return true
	}

	// Run command
	result := commandStage.runCommand()
	if result == false {
		log.Errorf("[command] exec: failed stage \"%s\"", commandStage.BaseStage.StageName)
	}
	return result
}

func (commandStage *CommandStage) waitFor() {
	if commandStage.WaitFor == "" {
		return
	}
	cond, _ := ParseWaitFor(commandStage.WaitFor) // TODO: error handling
	cond.Wait()
}

func (commandStage *CommandStage) runOnlyIf() bool {
	if commandStage.OnlyIf == "" {
		return true
	}
	log.Infof("[command] only_if: found \"only_if\" attribute in stage \"%s\"", commandStage.BaseStage.StageName)
	cmd := exec.Command("sh", "-c", commandStage.OnlyIf)
	log.Infof("[command] only_if: %s", commandStage.BaseStage.StageName)
	log.Debugf("[command] only_if literal: %s", commandStage.OnlyIf)
	cmd.Dir = commandStage.Directory
	result, _, _ := execCommand(cmd, "only_if", commandStage.BaseStage.StageName)
	return result
}

func (commandStage *CommandStage) runCommand() bool {
	cmd := exec.Command("sh", "-c", commandStage.Command)
	log.Infof("[command] exec: %s", commandStage.BaseStage.StageName)
	log.Debugf("[command] exec command literal: %s", commandStage.Command)
	cmd.Dir = commandStage.Directory
	commandStage.SetStart(time.Now().Unix())
	result, outResult, errResult := execCommand(cmd, "exec", commandStage.BaseStage.StageName)
	commandStage.SetEnd(time.Now().Unix())
	commandStage.SetOutResult(*outResult)
	commandStage.SetErrResult(*errResult)
	commandStage.SetReturnValue(result)
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

//AddCommand registers the specified command.
func (commandStage *CommandStage) AddCommand(command string) {
	commandStage.Command = command
	commandStage.BaseStage.Runner = commandStage
}

//SetDirectory sets the directory where the command is executed.
func (commandStage *CommandStage) SetDirectory(directory string) {
	commandStage.Directory = directory
}

//NewCommandStage creates one CommandStage object.
func NewCommandStage() *CommandStage {
	stage := CommandStage{Directory: "."}
	return &stage
}
