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

	"github.com/recruit-tech/plumber/log"
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
	var stage Stage
	switch stageType {
	case "command":
		stage = new(CommandStage)
	}

	prepareCh(&stage)
	return stage
}

func prepareCh(stage *Stage) {
	in := make(chan Mediator)
	out := make(chan Mediator)
	(*stage).SetInputCh(&in)
	(*stage).SetOutputCh(&out)
}

func ExecuteStage(stage Stage, monitorCh *chan Mediator) {
	mediator := <-*stage.GetInputCh()

	log.Debugf("mediator received: %v", mediator)
	log.Debugf("execute as parent: %v", stage)
	log.Debugf("execute as parent name %v", stage.GetStageName())

	result := stage.(Runner).Run()
	log.Debugf("execute as parent result %v", result)

	mediator.States[stage.GetStageName()] = fmt.Sprintf("%v", result)

	setChildStatus(&stage, &mediator, "waiting")

	if childStages := stage.GetChildStages(); childStages.Len() > 0 {
		log.Debugf("execute childstage: %v", childStages)

		for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
			log.Debugf("child name %+v\n", childStage.Value.(Stage).GetStageName())
			childInputCh := *childStage.Value.(Stage).GetInputCh()

			name := childStage.Value.(Stage).GetStageName()
			mediator.States[name] = fmt.Sprintf("%v", "waiting")

			go func(m Mediator) {
				childInputCh <- m
				close(childInputCh)
			}(mediator)

			go func(stage Stage) {
				ExecuteStage(stage, monitorCh)
			}(childStage.Value.(Stage))
		}

		for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
			go func(stage Stage) {
				childReceived := <-*stage.GetOutputCh()
				name := stage.GetStageName()
				log.Debugf("child state %+v\n", childReceived.States[name])
				mediator.States[name] = fmt.Sprintf("%v", childReceived.States[name])
			}(childStage.Value.(Stage))
		}
	}

	go func() {
		*stage.GetOutputCh() <- mediator
		close(*stage.GetOutputCh())
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
		log.Debugf("closing monitor channel.. %v\n", mediator)
		close(*monitCh)
	}
}

func Execute(stage Stage, mediator Mediator) Mediator {
	monitorCh := make(chan Mediator)
	name := stage.GetStageName()
	log.Debugf("----- Execute %v start ------\n", name)

	mediator.States[name] = fmt.Sprintf("%v", "waiting")

	go func(mediator Mediator) {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}(mediator)

	var lastReceive Mediator

	go ExecuteStage(stage, &monitorCh)

	for {
		receive, ok := <-monitorCh
		if !ok {
			log.Debugf("monitorCh closed")
			log.Debugf("----- Execute %v done ------\n\n", name)
			return lastReceive
		}
		log.Debugf("monitorCh received  %+v\n", receive)
		lastReceive = receive
	}
}
