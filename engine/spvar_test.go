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
	"testing"

	"github.com/recruit-tech/walter/pipelines"
	"github.com/recruit-tech/walter/stages"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	pipeline := pipelines.NewPipeline()
	pipeline.AddStage(new(stages.StageBuilder).NewStage("command").SetName("first").SetTarget("echo foobar").SetOutResult("foobar").Build())
	pipeline.AddStage(new(stages.StageBuilder).NewStage("command").SetName("second").SetTarget("echo baz").SetOutResult("baz").Build())
	spvar := NewSecialVariables(pipeline)
	result, err := spvar.Replace("__OUT[\"first\"] || __OUT[\"second\"]")
	assert.Nil(t, err)
	assert.Equal(t, "foobar || baz", result)
}
