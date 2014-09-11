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
	Type   string
}

func InitStage(stageType string) Stage {
	var stage Stage
	switch stageType {
	case "command":
		stage = new(CommandStage)
	}

	prepareCh(stage)
	return stage
}

func prepareCh(stage Stage) {
	in := make(chan Mediator)
	out := make(chan Mediator)
	stage.SetInputCh(&in)
	stage.SetOutputCh(&out)
}

func receiveInputs(inputCh *chan Mediator) []Mediator {
	mediatorsReceived := make([]Mediator, 0)
	for {
		m, ok := <-*inputCh
		if !ok {
			break
		}
		log.Debugf("received input: %+v", m)
		mediatorsReceived = append(mediatorsReceived, m)
	}
	return mediatorsReceived
}

func ExecuteStage(stage Stage, monitorCh *chan Mediator) {
	log.Debug("receiveing input")

	mediatorsReceived := receiveInputs(stage.GetInputCh())

	log.Debugf("received input size: %v", len(mediatorsReceived))
	log.Debugf("mediator received: %+v", mediatorsReceived)
	log.Debugf("execute as parent: %+v", stage)
	log.Debugf("execute as parent name %+v", stage.GetStageName())

	result := stage.(Runner).Run()
	log.Debugf("stage executution results: %+v, %+v", stage.GetStageName(), result)

	mediator := Mediator{States: make(map[string]string)}
	mediator.States[stage.GetStageName()] = fmt.Sprintf("%v", result)

	if childStages := stage.GetChildStages(); childStages.Len() > 0 {
		log.Debugf("execute childstage: %v", childStages)
		executeAllChildStages(&childStages, mediator, monitorCh)
		waitAllChildStages(&childStages, &stage)
	}

	log.Debugf("sending output of stage: %+v %v", stage.GetStageName(), mediator)
	*stage.GetOutputCh() <- mediator
	log.Debugf("closing output of stage: %+v", stage.GetStageName())
	close(*stage.GetOutputCh())

	for _, m := range mediatorsReceived {
		*monitorCh <- m
	}
	*monitorCh <- mediator

	finalizeMonitorChAfterExecute(mediatorsReceived, monitorCh)
}

func executeAllChildStages(childStages *list.List, mediator Mediator, monitorCh *chan Mediator) {
	for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
		log.Debugf("child name %+v\n", childStage.Value.(Stage).GetStageName())
		childInputCh := *childStage.Value.(Stage).GetInputCh()

		go func(stage Stage) {
			ExecuteStage(stage, monitorCh)
		}(childStage.Value.(Stage))

		log.Debugf("input child: %+v", mediator)
		childInputCh <- mediator
		log.Debugf("closing input: %+v", childStage.Value.(Stage).GetStageName())
		close(childInputCh)
	}
}

func waitAllChildStages(childStages *list.List, stage *Stage) {
	for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
		s := childStage.Value.(Stage)
		for {
			log.Debugf("receiving child: %v", s.GetStageName())
			childReceived, ok := <-*s.GetOutputCh()
			if !ok {
				log.Debug("closing child output")
				break
			}
			log.Debugf("sending child: %v", childReceived)
			*(*stage).GetOutputCh() <- childReceived
			log.Debugf("send child: %v", childReceived)
		}
		log.Debugf("finished executing child: %v", s.GetStageName())
	}
}

func finalizeMonitorChAfterExecute(mediators []Mediator, mon *chan Mediator) {
	if mediators[0].Type == "start" {
		log.Debug("finalize monitor channel..")
		mediatorEnd := Mediator{States: make(map[string]string), Type: "end"}
		*mon <- mediatorEnd
	} else {
		log.Debugf("skipped finalizing")
	}
}

func setChildStatus(stage *Stage, mediator *Mediator, status string) {
	if childStages := (*stage).GetChildStages(); childStages.Len() > 0 {
		for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
			name := childStage.Value.(Stage).GetStageName()
			mediator.States[name] = fmt.Sprintf("%v", status)
		}
	}
}

func Execute(stage Stage, mediator Mediator) Mediator {
	monitorCh := make(chan Mediator)
	mediator.Type = "start"
	name := stage.GetStageName()
	log.Debugf("----- Execute %v start ------\n", name)

	go func(mediator Mediator) {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}(mediator)

	go ExecuteStage(stage, &monitorCh)

	for {
		receive, ok := <-*stage.GetOutputCh()
		if !ok {
			log.Debugf("outputCh closed")
			break
		}
		log.Debugf("outputCh received  %+v\n", receive)
	}

	for {
		receive := <-monitorCh
		if receive.Type == "end" {
			log.Debugf("monitorCh closed")
			log.Debugf("monitorCh last received:  %+v\n", receive)
			log.Debugf("----- Execute %v done ------\n\n", name)
			return receive
		}
		log.Debugf("monitorCh received  %+v\n", receive)
	}
}
