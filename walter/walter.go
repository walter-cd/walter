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
package walter

import (
	"fmt"

	"github.com/recruit-tech/walter/config"
	"github.com/recruit-tech/walter/engine"
	"github.com/recruit-tech/walter/log"
	"github.com/recruit-tech/walter/service"
	"github.com/recruit-tech/walter/stages"
)

type Walter struct {
	Engine *engine.Engine
}

func New(opts *config.Opts) (*Walter, error) {
	configData := config.ReadConfig(opts.PipelineFilePath)
	pipeline, err := config.Parse(configData)
	if err != nil {
		return nil, err
	}
	monitorCh := make(chan stages.Mediator)
	engine := &engine.Engine{
		Pipeline:  pipeline,
		Opts:      opts,
		MonitorCh: &monitorCh,
	}
	return &Walter{
		Engine: engine,
	}, err
}

func (e *Walter) Run() bool {
	if e.Engine.Opts.Mode == "local" {
		mediator := e.Engine.RunOnce()
		return !mediator.IsAnyFailure()
	} else {
		// load .walter-update
		log.Info("loading update file...")
		update, err := service.LoadLastUpdate(e.Engine.Pipeline.RepoService.GetUpdateFilePath())
		if err != nil {
			log.Warnf("failed to load update: %s", err)
		}

		// get latest commti and pull requests
		log.Info("downloading commits and pull requests...")
		commits, err := e.Engine.Pipeline.RepoService.GetCommits(update)
		if err != nil {
			log.Errorf("failed to get commits: %s", err)
			return false
		}

		log.Info("suceeded to get commits")
		for e := commits.Front(); e != nil; e = e.Next() {
			fmt.Println(e) // TODO implement Running walter with local mode
		}

		// save .walter-update
		log.Info("saving update file...")
		result := service.SaveLastUpdate(e.Engine.Pipeline.RepoService.GetUpdateFilePath(), update)
		if result == false {
			log.Warnf("failed to save update")
		}
		return true
	}
}
