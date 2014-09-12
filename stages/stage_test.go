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

	stage.AddCommand("ls")
	mon := make(chan Mediator)

	go ExecuteStage(stage, &mon)

	mediator := Mediator{States: make(map[string]string), Type: "start"}
	go func() {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}()

	for {
		_, ok := <-*stage.GetOutputCh()
		if !ok {
			break
		}
	}

	var m Mediator
	var actual Mediator
	for {
		m = <-mon
		if m.Type == "end" {
			break
		}
		actual = m
	}

	expected := "true"

	if expected != actual.States[stage.StageName] {
		t.Errorf("got %v\nwant %v", actual.States[stage.StageName], expected)
	}
}

func TestExecuteWithSingleStageFailed(t *testing.T) {
	stage := &CommandStage{}
	stage.StageName = "test_command_stage"
	prepareCh(stage)

	stage.AddCommand("nothingcommand")
	mon := make(chan Mediator)

	go ExecuteStage(stage, &mon)

	mediator := Mediator{States: make(map[string]string), Type: "start"}
	go func() {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}()

	for {
		_, ok := <-*stage.GetOutputCh()
		if !ok {
			break
		}
	}

	var m Mediator
	var actual Mediator
	for {
		m = <-mon
		if m.Type == "end" {
			break
		}
		actual = m
	}

	expected := "false"

	if expected != actual.States[stage.StageName] {
		t.Errorf("got %v\nwant %v", actual.States[stage.StageName], expected)
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
	mediator := Mediator{States: make(map[string]string)}
	mediator.Type = "start"

	t.Logf("execute: %+v", stage)

	go func() {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}()

	go ExecuteStage(stage, &mon)

	for {
		_, ok := <-*stage.GetOutputCh()
		if !ok {
			break
		}
	}

	var m Mediator
	var actual Mediator
	for {
		m = <-mon
		if m.Type == "end" {
			break
		}
		actual = m
	}

	expected := "true"

	if expected != actual.States[stage.StageName] {
		t.Errorf("got %v\nwant %v", actual.States[stage.StageName], expected)
	}
}

func TestExecuteWithSingleStageHasErrChild(t *testing.T) {
	stage := &CommandStage{}
	stage.StageName = "test_command_stage"
	prepareCh(stage)

	child := &CommandStage{}
	child.StageName = "test_child"
	prepareCh(child)

	stage.AddCommand("ls -l")
	child.AddCommand("nothingcommand")

	stage.AddChildStage(child)

	mon := make(chan Mediator)
	mediator := Mediator{States: make(map[string]string)}
	mediator.Type = "start"

	t.Logf("execute: %+v", stage)

	go func() {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}()

	go ExecuteStage(stage, &mon)

	for {
		_, ok := <-*stage.GetOutputCh()
		if !ok {
			break
		}
	}

	var m Mediator
	acm := Mediator{States: make(map[string]string)}
	for {
		m = <-mon
		for k, v := range m.States {
			acm.States[k] = v
		}
		if m.Type == "end" {
			break
		}
	}

	t.Logf("accumulated output: %+v", acm)
	expected := "true"

	if expected != acm.States[stage.StageName] {
		t.Errorf("got %v\nwant %v", acm.States[stage.StageName], expected)
	}

	expected = "false"
	if expected != acm.States[child.StageName] {
		t.Errorf("got %v\nwant %v", acm.States[child.StageName], expected)
	}
}
