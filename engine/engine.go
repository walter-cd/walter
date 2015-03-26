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
package engine

import (
	"container/list"
	"fmt"

	"github.com/recruit-tech/walter/config"
	"github.com/recruit-tech/walter/log"
	"github.com/recruit-tech/walter/pipelines"
	"github.com/recruit-tech/walter/stages"
)

type Engine struct {
	Resources *pipelines.Resources
	MonitorCh *chan stages.Mediator
	Opts      *config.Opts
}

type Result struct {
	Pipeline *stages.Mediator
	Cleanup  *stages.Mediator
}

func (r *Result) IsSucceeded() bool {
	if !r.Pipeline.IsAnyFailure() && !r.Cleanup.IsAnyFailure() {
		return true
	}
	return false
}

func (e *Engine) RunOnce() *Result {
	pipe_result := e.executePipeline(e.Resources.Pipeline, "pipeline")
	cleanup_result := e.executePipeline(e.Resources.Cleanup, "cleanup")
	return &Result{Pipeline: &pipe_result, Cleanup: &cleanup_result}
}

func (e *Engine) executePipeline(pipeline *pipelines.Pipeline, name string) stages.Mediator {
	log.Infof("Preparing to run %s process...", name)
	var mediator stages.Mediator
	for stageItem := pipeline.Stages.Front(); stageItem != nil; stageItem = stageItem.Next() {
		log.Debugf("Executing planned stage: %s\n", stageItem.Value)
		mediator = e.Execute(stageItem.Value.(stages.Stage), mediator)
	}
	log.Infof("Finished running %s process...", name)
	return mediator
}

func (e *Engine) receiveInputs(inputCh *chan stages.Mediator) []stages.Mediator {
	mediatorsReceived := make([]stages.Mediator, 0)
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

func (e *Engine) ExecuteStage(stage stages.Stage) {
	log.Debug("Receiving input")

	mediatorsReceived := e.receiveInputs(stage.GetInputCh())

	log.Debugf("Received input size: %v", len(mediatorsReceived))
	log.Debugf("Mediator received: %+v", mediatorsReceived)
	log.Debugf("Execute as parent: %+v", stage)
	log.Debugf("Execute as parent name %+v", stage.GetStageName())

	var result bool
	if !e.isUpstreamAnyFailure(mediatorsReceived) || e.Opts.StopOnAnyFailure {
		result = stage.(stages.Runner).Run()
	} else {
		log.Warnf("Execution is skipped: %v", stage.GetStageName())
		result = false
	}
	log.Debugf("Stage execution results: %+v, %+v", stage.GetStageName(), result)
	e.Resources.ReportStageResult(stage, result)

	mediator := stages.Mediator{States: make(map[string]string)}
	mediator.States[stage.GetStageName()] = fmt.Sprintf("%v", result)

	if childStages := stage.GetChildStages(); childStages.Len() > 0 {
		log.Debugf("Execute childstage: %v", childStages)
		e.executeAllChildStages(&childStages, mediator)
		e.waitAllChildStages(&childStages, &stage)
	}

	log.Debugf("Sending output of stage: %+v %v", stage.GetStageName(), mediator)
	*stage.GetOutputCh() <- mediator
	log.Debugf("Closing output of stage: %+v", stage.GetStageName())
	close(*stage.GetOutputCh())

	for _, m := range mediatorsReceived {
		*e.MonitorCh <- m
	}
	*e.MonitorCh <- mediator

	e.finalizeMonitorChAfterExecute(mediatorsReceived)
}

func (e *Engine) isUpstreamAnyFailure(mediators []stages.Mediator) bool {
	for _, m := range mediators {
		if m.IsAnyFailure() == true {
			return true
		}
	}
	return false
}

func (e *Engine) executeAllChildStages(childStages *list.List, mediator stages.Mediator) {
	for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
		log.Debugf("Child name %+v\n", childStage.Value.(stages.Stage).GetStageName())
		childInputCh := *childStage.Value.(stages.Stage).GetInputCh()

		go func(stage stages.Stage) {
			e.ExecuteStage(stage)
		}(childStage.Value.(stages.Stage))

		log.Debugf("Input child: %+v", mediator)
		childInputCh <- mediator
		log.Debugf("Closing input: %+v", childStage.Value.(stages.Stage).GetStageName())
		close(childInputCh)
	}
}

func (e *Engine) waitAllChildStages(childStages *list.List, stage *stages.Stage) {
	for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
		s := childStage.Value.(stages.Stage)
		for {
			log.Debugf("Receiving child: %v", s.GetStageName())
			childReceived, ok := <-*s.GetOutputCh()
			if !ok {
				log.Debug("Closing child output")
				break
			}
			log.Debugf("Sending child: %v", childReceived)
			*(*stage).GetOutputCh() <- childReceived
			log.Debugf("Send child: %v", childReceived)
		}
		log.Debugf("Finished executing child: %v", s.GetStageName())
	}
}

func (e *Engine) finalizeMonitorChAfterExecute(mediators []stages.Mediator) {
	if mediators[0].Type == "start" {
		log.Debug("Finalize monitor channel..")
		mediatorEnd := stages.Mediator{States: make(map[string]string), Type: "end"}
		*e.MonitorCh <- mediatorEnd
	} else {
		log.Debugf("Skipped finalizing")
	}
}

func (e *Engine) Execute(stage stages.Stage, mediator stages.Mediator) stages.Mediator {
	monitorCh := e.MonitorCh
	mediator.Type = "start"
	name := stage.GetStageName()
	log.Debugf("----- Execute %v start ------\n", name)

	go func(mediator stages.Mediator) {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}(mediator)

	go e.ExecuteStage(stage)

	for {
		receive, ok := <-*stage.GetOutputCh()
		if !ok {
			log.Debugf("outputCh closed")
			break
		}
		log.Debugf("outputCh received  %+v\n", receive)
	}

	receives := make([]stages.Mediator, 0)
	for {
		receive := <-*monitorCh
		receives = append(receives, receive)
		if receive.Type == "end" {
			log.Debugf("monitorCh closed")
			log.Debugf("monitorCh last received:  %+v\n", receive)
			log.Debugf("----- Execute %v done ------\n\n", name)
			return e.bindReceives(&receives)
		}
		log.Debugf("monitorCh received  %+v\n", receive)
	}
}

func (e *Engine) bindReceives(rs *[]stages.Mediator) stages.Mediator {
	ret := &stages.Mediator{States: make(map[string]string)}
	for _, r := range *rs {
		for k, v := range r.States {
			ret.States[k] = v
		}
	}
	return *ret
}
