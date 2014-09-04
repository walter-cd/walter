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
	"container/list"
	"fmt"
)

type Stage interface {
	AddChildStage(Stage)
	GetChildStages() list.List
	GetStageName() string
	SetStageName(string)
	GetInputCh() *chan Mediator
	SetInputCh(*chan Mediator)
	GetOutputCh() *chan Mediator
	SetOutputCh(*chan Mediator)
}

type Runner interface {
	Run() bool
}

type Mediator struct {
	States map[string]string
}

type BaseStage struct {
	Runner
	InputCh     *chan Mediator
	OutputCh    *chan Mediator
	ChildStages list.List
	StageName   string `config:"stage_name"`
}

func InitStage(stageType string) Stage {
	switch stageType {
	case "command":
		return new(CommandStage)
	}
	return nil
}

func (b *BaseStage) Run() bool {
	fmt.Println("called BaseStage.Run")
	if b.Runner == nil {
		panic("Mast have a child class assigned")
	}
	return b.Runner.Run()
}

func (b *BaseStage) AddChildStage(stage Stage) {
	fmt.Println("added childstage: %v", stage)
	b.ChildStages.PushBack(stage)
}

func (b *BaseStage) GetChildStages() list.List {
	return b.ChildStages
}

func (b *BaseStage) GetStageName() string {
	return b.StageName
}

func (b *BaseStage) SetStageName(stageName string) {
	b.StageName = stageName
}

func (b *BaseStage) SetInputCh(inputCh *chan Mediator) {
	b.InputCh = inputCh
}

func (b *BaseStage) GetInputCh() *chan Mediator {
	return b.InputCh
}

func (b *BaseStage) SetOutputCh(outputCh *chan Mediator) {
	b.OutputCh = outputCh
}

func (b *BaseStage) GetOutputCh() *chan Mediator {
	return b.OutputCh
}

func ExecuteStage(stage Stage, inputChan *chan Mediator, outputChan *chan Mediator, monitorChan *chan Mediator) {
	mediator := <-*inputChan

	fmt.Println("mediator received: %v", mediator)
	fmt.Println("execute as parent: %v", stage)
	fmt.Println("execute as parent name %v", stage.GetStageName())

	result := stage.(Runner).Run()
	fmt.Println("execute as parent result %v", result)

	mediator.States[stage.GetStageName()] = fmt.Sprintf("%v", result)

	setChildStatus(&stage, &mediator, "waiting")

	if childStages := stage.GetChildStages(); childStages.Len() > 0 {
		fmt.Println("execute childstage: %v", childStages)

		for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
			fmt.Printf("child name %+v\n", childStage.Value.(Stage).GetStageName())
			childInputCh := make(chan Mediator)
			childOutputCh := make(chan Mediator)
			childStage.Value.(Stage).SetInputCh(&childInputCh)
			childStage.Value.(Stage).SetOutputCh(&childOutputCh)

			name := childStage.Value.(Stage).GetStageName()
			mediator.States[name] = fmt.Sprintf("%v", "waiting")

			go func(m Mediator) {
				childInputCh <- m
				close(childInputCh)
			}(mediator)

			go func(stage Stage) {
				ExecuteStage(stage, &childInputCh, &childOutputCh, monitorChan)
			}(childStage.Value.(Stage))
		}

		for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
			go func(stage Stage) {
				childReceived := <-*stage.GetOutputCh()
				name := stage.GetStageName()
				fmt.Printf("child state %+v\n", childReceived.States[name])
				mediator.States[name] = fmt.Sprintf("%v", childReceived.States[name])
			}(childStage.Value.(Stage))
		}
	}

	go func() {
		*outputChan <- mediator
		close(*outputChan)
	}()
	*monitorChan <- mediator

	closeAfterExecute(&mediator, monitorChan)
}

func setChildStatus(stage *Stage, mediator *Mediator, status string) {
	if childStages := (*stage).GetChildStages(); childStages.Len() > 0 {
		for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
			name := childStage.Value.(Stage).GetStageName()
			mediator.States[name] = fmt.Sprintf("%v", status)
		}
	}
}

func closeAfterExecute(mediator *Mediator, monitCh *chan Mediator) {
	allDone := true
	for _, v := range mediator.States {
		if v == "waiting" {
			allDone = false
		}
	}

	if allDone {
		fmt.Printf("closing monitor channel.. %v\n", mediator)
		close(*monitCh)
	}
}

func Execute(stage Stage, mediator Mediator) Mediator {
	inputChan := make(chan Mediator)
	outputChan := make(chan Mediator)
	monitorChan := make(chan Mediator)

	name := stage.GetStageName()
	fmt.Printf("----- Execute %v start ------\n", name)

	mediator.States[name] = fmt.Sprintf("%v", "waiting")

	go func(mediator Mediator) {
		inputChan <- mediator
		outputChan <- mediator
		close(outputChan)
		close(inputChan)
	}(mediator)

	var lastReceive Mediator

	go ExecuteStage(stage, &inputChan, &outputChan, &monitorChan)

	for {
		receive, ok := <-monitorChan
		if !ok {
			fmt.Println("monitorChan closed")
			fmt.Printf("----- Execute %v done ------\n\n", name)
			return lastReceive
		}
		fmt.Printf("monitorChan received  %+v\n", receive)
		lastReceive = receive
	}
}
