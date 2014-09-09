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
	"testing"
)

func TestExecuteWithSingleStage(t *testing.T) {
	stage := &CommandStage{}
	stage.StageName = "test_command_stage"
	prepareCh(stage)

	stage.AddCommand("ls -l")
	mon := make(chan Mediator)

	go ExecuteStage(stage, &mon)

	mediator := Mediator{States: make(map[string]string)}
	*stage.GetInputCh() <- mediator
	m := <-mon

	actual := m.States[stage.StageName]
	expected := "true"

	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestExecuteWithSingleStageFailed(t *testing.T) {
	stage := &CommandStage{}
	stage.StageName = "test_command_stage"
	prepareCh(stage)

	stage.AddCommand("nothingcommand")
	mon := make(chan Mediator)

	go ExecuteStage(stage, &mon)

	mediator := Mediator{States: make(map[string]string)}
	*stage.GetInputCh() <- mediator
	m := <-mon

	actual := m.States[stage.StageName]
	expected := "false"

	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestExecuteWithSingleStageHasChild(t *testing.T) {
	stage := &CommandStage{}
	stage.StageName = "test_command_stage"
	prepareCh(stage)

	child := &CommandStage{}
	child.StageName = "test_child"
	prepareCh(child)

	stage.AddCommand("ls -l")
	child.AddCommand("ls -l")

	stage.AddChildStage(child)

	mon := make(chan Mediator)

	go ExecuteStage(stage, &mon)

	mediator := Mediator{States: make(map[string]string)}
	*stage.GetInputCh() <- mediator

	var m Mediator
	for i := 0; i < 2; i++ {
		m = <-mon
	}

	actual := m.States[stage.StageName]
	expected := "true"

	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
