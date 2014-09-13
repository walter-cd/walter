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
package engine

import (
	"testing"

	"github.com/recruit-tech/plumber/config"
	"github.com/recruit-tech/plumber/pipelines"
	"github.com/recruit-tech/plumber/stages"
)

func createCommandStage(command string) *stages.CommandStage {
	in := make(chan stages.Mediator)
	out := make(chan stages.Mediator)
	return &stages.CommandStage{
		Command: command,
		BaseStage: stages.BaseStage{
			StageName: command,
			InputCh:   &in,
			OutputCh:  &out,
		},
	}
}

func TestRunOnce(t *testing.T) {
	pipeline := pipelines.NewPipeline()
	pipeline.AddStage(createCommandStage("echo foobar"))
	pipeline.AddStage(createCommandStage("echo baz"))
	monitorCh := make(chan stages.Mediator)
	engine := &Engine{
		Pipeline:  pipeline,
		MonitorCh: &monitorCh,
	}
	m := engine.RunOnce()

	expected := "true"
	actual := m.States["echo foobar"]
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestRunOnceWithOptsOffStopOnAnyFailure(t *testing.T) {
	pipeline := pipelines.NewPipeline()
	pipeline.AddStage(createCommandStage("echo foobar"))
	pipeline.AddStage(createCommandStage("thisiserrorcommand"))
	pipeline.AddStage(createCommandStage("echo foobar2"))
	monitorCh := make(chan stages.Mediator)
	o := &config.Opts{StopOnAnyFailure: false}
	engine := &Engine{
		Pipeline:  pipeline,
		MonitorCh: &monitorCh,
		Opts:      o,
	}

	m := engine.RunOnce()

	expected := "false"
	actual := m.States["echo foobar2"]

	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestRunOnceWithOptsOnStopOnAnyFailure(t *testing.T) {
	pipeline := pipelines.NewPipeline()
	pipeline.AddStage(createCommandStage("echo foobar"))
	pipeline.AddStage(createCommandStage("thisiserrorcommand"))
	pipeline.AddStage(createCommandStage("echo foobar2"))
	monitorCh := make(chan stages.Mediator)
	o := &config.Opts{StopOnAnyFailure: true}
	engine := &Engine{
		Pipeline:  pipeline,
		MonitorCh: &monitorCh,
		Opts:      o,
	}

	m := engine.RunOnce()

	expected := "true"
	actual := m.States["echo foobar2"]

	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestExecuteWithSingleStage(t *testing.T) {
	stage := &stages.CommandStage{}
	stage.StageName = "test_command_stage"
	stages.PrepareCh(stage)

	stage.AddCommand("ls")
	mon := make(chan stages.Mediator)
	e := &Engine{
		MonitorCh: &mon,
	}

	go e.ExecuteStage(stage)

	mediator := stages.Mediator{States: make(map[string]string), Type: "start"}
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

	var m stages.Mediator
	var actual stages.Mediator
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
	stage := &stages.CommandStage{}
	stage.StageName = "test_command_stage"
	stages.PrepareCh(stage)

	stage.AddCommand("nothingcommand")
	mon := make(chan stages.Mediator)
	e := &Engine{
		MonitorCh: &mon,
	}

	go e.ExecuteStage(stage)

	mediator := stages.Mediator{States: make(map[string]string), Type: "start"}
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

	var m stages.Mediator
	var actual stages.Mediator
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
	stage := &stages.CommandStage{}
	stage.StageName = "test_command_stage"
	stages.PrepareCh(stage)

	child := &stages.CommandStage{}
	child.StageName = "test_child"
	stages.PrepareCh(child)

	stage.AddCommand("ls -l")
	child.AddCommand("ls -l")

	stage.AddChildStage(child)

	mon := make(chan stages.Mediator)
	mediator := stages.Mediator{States: make(map[string]string)}
	mediator.Type = "start"

	e := &Engine{
		MonitorCh: &mon,
	}

	t.Logf("execute: %+v", stage)

	go func() {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}()

	go e.ExecuteStage(stage)

	for {
		_, ok := <-*stage.GetOutputCh()
		if !ok {
			break
		}
	}

	var m stages.Mediator
	var actual stages.Mediator
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
	stage := &stages.CommandStage{}
	stage.StageName = "test_command_stage"
	stages.PrepareCh(stage)

	child := &stages.CommandStage{}
	child.StageName = "test_child"
	stages.PrepareCh(child)

	stage.AddCommand("ls -l")
	child.AddCommand("nothingcommand")

	stage.AddChildStage(child)

	mon := make(chan stages.Mediator)
	mediator := stages.Mediator{States: make(map[string]string)}
	mediator.Type = "start"

	e := &Engine{
		MonitorCh: &mon,
	}

	t.Logf("execute: %+v", stage)

	go func() {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}()

	go e.ExecuteStage(stage)

	for {
		_, ok := <-*stage.GetOutputCh()
		if !ok {
			break
		}
	}

	var m stages.Mediator
	acm := stages.Mediator{States: make(map[string]string)}
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
