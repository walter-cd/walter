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

func InitStage(stageType string) Stage {
	switch stageType {
	case "command":
		return new(CommandStage)
	}
	return nil
}

func ExecuteStage(stage Stage, inputCh *chan Mediator, outputCh *chan Mediator, monitorCh *chan Mediator) {
	mediator := <-*inputCh

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
				ExecuteStage(stage, &childInputCh, &childOutputCh, monitorCh)
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
		*outputCh <- mediator
		close(*outputCh)
	}()
	*monitorCh <- mediator

	closeAfterExecute(&mediator, monitorCh)
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
	inputCh := make(chan Mediator)
	outputCh := make(chan Mediator)
	monitorCh := make(chan Mediator)

	name := stage.GetStageName()
	fmt.Printf("----- Execute %v start ------\n", name)

	mediator.States[name] = fmt.Sprintf("%v", "waiting")

	go func(mediator Mediator) {
		inputCh <- mediator
		close(inputCh)
	}(mediator)

	var lastReceive Mediator

	go ExecuteStage(stage, &inputCh, &outputCh, &monitorCh)

	for {
		receive, ok := <-monitorCh
		if !ok {
			fmt.Println("monitorCh closed")
			fmt.Printf("----- Execute %v done ------\n\n", name)
			return lastReceive
		}
		fmt.Printf("monitorCh received  %+v\n", receive)
		lastReceive = receive
	}
}
