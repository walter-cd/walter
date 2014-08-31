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

func (b *BaseStage) setInputCh() {
	inputCh := make(chan Mediator)
	b.InputCh = &inputCh
}

func (b *BaseStage) getInputCh() *chan Mediator {
	return b.InputCh
}

func ExecuteStage(stage Stage, outputChan *chan Mediator) chan Mediator {
	receiver := make(chan Mediator)
	var mediator Mediator
	for {
		received, ok := <-*outputChan
		if !ok {
			break
		}
		mediator = received
	}
	fmt.Println("mediator received: %v", mediator)
	fmt.Println("execute parent: %v", stage)

	go func(stage Stage) {
		fmt.Println("execute parent name %v", stage.GetStageName())
		result := stage.(Runner).Run()
		fmt.Println("execute parent result %v", result)
		mediator.States[stage.GetStageName()] = fmt.Sprintf("%v", result)
		receiver <- mediator
		if childStages := stage.GetChildStages(); childStages.Len() > 0 {
			fmt.Println("execute childstage: %v", childStages)
			for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
				//ExecuteStage(childStage.Value.(Stage))
			}
		}
		close(receiver)
	}(stage)

	return receiver
}

func Execute(stage Stage, mediator Mediator) Mediator {
	outputChan := make(chan Mediator)

	go func(mediator Mediator) {
		outputChan <- mediator
		close(outputChan)
	}(mediator)

	receiver := ExecuteStage(stage, &outputChan)

	var lastReceive Mediator
	for {
		receive, ok := <-receiver
		if !ok {
			fmt.Println("receiver closed")
			return lastReceive
		}
		lastReceive = receive
	}
}
