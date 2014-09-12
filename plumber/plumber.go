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
package plumber

import (
	"github.com/recruit-tech/plumber/config"
	"github.com/recruit-tech/plumber/engine"
	"github.com/recruit-tech/plumber/pipelines"
	"github.com/recruit-tech/plumber/stages"
)

type Plumber struct {
	Pipeline *pipelines.Pipeline
	Engine   *engine.Engine
}

func New(opts *config.Opts) *Plumber {
	configData := config.ReadConfig(opts.PipelineFilePath)
	pipeline := (config.Parse(configData))
	monitorCh := make(chan stages.Mediator)
	engine := &engine.Engine{
		Pipeline:  pipeline,
		Opts:      opts,
		MonitorCh: &monitorCh,
	}
	return &Plumber{
		Engine: engine,
	}
}

func (e *Plumber) Run() {
	e.Engine.RunOnce()
}
