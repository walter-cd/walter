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

// Package engine defines the Engine struct, which execute registered pipelines.
package engine

import (
	"container/list"
	//	"os"
	"strconv"

	"github.com/walter-cd/walter/config"
	"github.com/walter-cd/walter/log"
	"github.com/walter-cd/walter/pipelines"
	"github.com/walter-cd/walter/stages"
)

// Engine executes the its pipeline.
type Engine struct {
	Resources    *pipelines.Resources
	MonitorCh    *chan stages.Mediator
	Opts         *config.Opts
	EnvVariables *config.EnvVariables
}

// Result stores the output in pipelines.
type Result struct {
	Pipeline *stages.Mediator
	Cleanup  *stages.Mediator
}

// IsSucceeded shows the pipeline finished successfully or not.
func (r *Result) IsSucceeded() bool {
	if !r.Pipeline.IsAnyFailure() && !r.Cleanup.IsAnyFailure() {
		return true
	}
	return false
}

// RunOnce executes the pipeline and the cleanup prccedures.
func (e *Engine) RunOnce() *Result {
	pipeResult := e.executePipeline(e.Resources.Pipeline, "pipeline")
	cleanupResult := e.executePipeline(e.Resources.Cleanup, "cleanup")
	return &Result{Pipeline: &pipeResult, Cleanup: &cleanupResult}
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
	var mediatorsReceived = make([]stages.Mediator, 0)
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

// ExecuteStage executes the supplied stage
func (e *Engine) ExecuteStage(stage stages.Stage) {
	log.Debug("Receiving input")
	mediatorsReceived := e.receiveInputs(stage.GetInputCh())

	log.Debugf("Received input size: %v", len(mediatorsReceived))
	log.Debugf("Mediator received: %+v", mediatorsReceived)
	log.Debugf("Execute as parent: %+v", stage)
	log.Debugf("Execute as parent name %+v", stage.GetStageName())

	mediator := stages.Mediator{States: make(map[string]string)}
	mediator.States[stage.GetStageName()] = e.executeStage(stage, mediatorsReceived, mediator)

	log.Debugf("Sending output of stage: %+v %v", stage.GetStageName(), mediator)
	*stage.GetOutputCh() <- mediator
	close(*stage.GetOutputCh())
	log.Debugf("Closed output of stage: %+v", stage.GetStageName())

	e.finalizeMonitorChAfterExecute(mediatorsReceived, mediator)
}

func (e *Engine) executeChildStages(stage *stages.Stage, mediator *stages.Mediator) {
	if childStages := (*stage).GetChildStages(); childStages.Len() > 0 {
		log.Debugf("Execute childstage: %v", childStages)
		e.executeAllChildStages(&childStages, *mediator)
		e.waitAllChildStages(&childStages, stage)
	}
}

func (e *Engine) executeStage(stage stages.Stage, received []stages.Mediator, mediator stages.Mediator) string {
	var result string
	if !e.isUpstreamAnyFailure(received) || e.Opts.StopOnAnyFailure {
		result = strconv.FormatBool(stage.(stages.Runner).Run())
		e.EnvVariables.ExportSpecialVariable("__OUT[\""+stage.GetStageName()+"\"]", stage.GetOutResult())
		e.EnvVariables.ExportSpecialVariable("__ERR[\""+stage.GetStageName()+"\"]", stage.GetErrResult())
		e.EnvVariables.ExportSpecialVariable("__COMBINED[\""+stage.GetStageName()+"\"]", stage.GetCombinedResult())
		e.EnvVariables.ExportSpecialVariable("__RESULT[\""+stage.GetStageName()+"\"]", result)
		e.executeChildStages(&stage, &mediator)
	} else {
		log.Warnf("Execution is skipped: %v", stage.GetStageName())
		if childStages := stage.GetChildStages(); childStages.Len() > 0 {
			for childStage := childStages.Front(); childStage != nil; childStage = childStage.Next() {
				log.Warnf("Execution of child stage is skipped: %v", childStage.Value.(stages.Stage).GetStageName())
			}
		}
		result = "skipped"
	}
	e.Resources.ReportStageResult(stage, result)
	log.Debugf("Stage execution results: %+v, %+v", stage.GetStageName(), result)
	return result
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

		go func() {
			e.ExecuteStage(childStage.Value.(stages.Stage))
		}()

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

func (e *Engine) finalizeMonitorChAfterExecute(mediators []stages.Mediator, mediator stages.Mediator) {
	// Append to monitorCh
	*e.MonitorCh <- mediator
	for _, m := range mediators {
		*e.MonitorCh <- m
	}

	// Set type
	if mediators[0].Type == "start" {
		log.Debug("Finalize monitor channel..")
		mediatorEnd := stages.Mediator{States: make(map[string]string), Type: "end"}
		*e.MonitorCh <- mediatorEnd
	} else {
		log.Debugf("Skipped finalizing")
	}
}

//Execute executes a stage using the supplied mediator
func (e *Engine) Execute(stage stages.Stage, mediator stages.Mediator) stages.Mediator {
	mediator.Type = "start"
	name := stage.GetStageName()
	log.Debugf("----- Execute %v start ------\n", name)

	go func() {
		*stage.GetInputCh() <- mediator
		close(*stage.GetInputCh())
	}()

	go e.ExecuteStage(stage)
	e.waitCloseOutputCh(stage)
	return e.waitMonitorChFinalized(name)
}

func (e *Engine) waitCloseOutputCh(stage stages.Stage) {
	for {
		receive, ok := <-*stage.GetOutputCh()
		if !ok {
			log.Debugf("outputCh closed")
			break
		}
		log.Debugf("outputCh received  %+v\n", receive)
	}
}

func (e *Engine) waitMonitorChFinalized(name string) stages.Mediator {
	var receives = make([]stages.Mediator, 0)
	for {
		receive := <-*e.MonitorCh
		receives = append(receives, receive)
		if receive.Type == "end" {
			log.Debugf("monitorCh closed")
			log.Debugf("monitorCh last received:  %+v\n", receive)
			log.Debugf("----- Execute %v done ------\n\n", name)
			return e.bindReceives(&receives)
		}
		log.Debugf("monitorCh received %+v\n", receive)
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
