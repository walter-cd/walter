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
	"net"
	"os"
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

// WaitFor wait until the condtions are satisfied
func (waitFor *WaitFor) Wait() {
	// delay
	if waitFor.Delay > 0.0 {
		log.Info("Wait specified time: " + strconv.FormatFloat(waitFor.Delay, 'f', 6, 64))
		time.Sleep(time.Duration(waitFor.Delay) * time.Second)
		return
	}

	// file created
	if waitFor.File != "" && (waitFor.State == "present" || waitFor.State == "ready") {
		log.Info("Wait for file: " + waitFor.File + " is created...")
		for {
			if isFileExist(waitFor.File) {
				log.Info("File: " + waitFor.File + " found.")
				return
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	// file removed
	if waitFor.File != "" && (waitFor.State == "absent" || waitFor.State == "unready") {
		log.Info("Wait for file: " + waitFor.File + " is removed...")
		for {
			if !isFileExist(waitFor.File) {
				log.Info("File: " + waitFor.File + " removed.")
				return
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	// port open
	if waitFor.Host != "" && waitFor.Port > 0 && (waitFor.State == "present" || waitFor.State == "ready") {
		log.Info("Wait for port: " + waitFor.Host + ":" + strconv.Itoa(waitFor.Port) + " is opened...")
		for {
			if isConnect(waitFor.Host, waitFor.Port) {
				return
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	// port close
	if waitFor.Host != "" && waitFor.Port > 0 && (waitFor.State == "absent" || waitFor.State == "unready") {
		log.Info("Wait for: " + waitFor.Host + ":" + strconv.Itoa(waitFor.Port) + " is closed...")

		for {
			if !isConnect(waitFor.Host, waitFor.Port) {
				return
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func isFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func isConnect(host string, port int) bool {
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// ParseWaitFor returns the WaitFor instance from given string
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

	if validateWaitForCondition(wait) {
		return wait, nil
	} else {
		return nil, errors.New("Illegal condition were found")
	}
}

func validateWaitForCondition(wait *WaitFor) bool {
	// check duplicate targets
	if wait.Port != 0 && wait.File != "" {
		log.Error("[command] Port and File are not able to be specified at the same time.")
		return false
	} else if wait.Port != 0 && wait.Delay > 0.0 {
		log.Error("[command] Port and Delay are not able to be specified at the same time.")
		return false
	} else if wait.File != "" && wait.Delay > 0.0 {
		log.Error("[command] File and Delay are not able to be specified at the same time.")
		return false
	}

	// check illegal conditions
	if wait.Delay < 0 {
		log.Error("[command] Delay must be positive.")
		return false
	} else if wait.Port < 0 {
		log.Error("[command] Port must be positive.")
		return false
	} else if wait.Port > 0 && wait.Host == "" {
		log.Error("[command] Host must be specified when port number is specified.")
		return false
	}

	// check illegal states
	if wait.State != "present" && wait.State != "ready" && wait.State != "absent" && wait.State != "unready" {
		log.Errorf("[command] \"%s\" is an unsupported state", wait.State)
		return false
	}

	// misc checks
	if wait.Port > 0 && wait.State == "" {
		log.Error("[command] State must be specified for port.")
		return false
	} else if wait.File != "" && wait.State == "" {
		log.Error("[command] State must be specified for file.")
		return false
	}

	return true
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
