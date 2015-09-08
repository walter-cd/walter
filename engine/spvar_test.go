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
	// "github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	pipeline := pipelines.NewPipeline()
	pipeline.AddStage(createCommandStageWithName("first", "ls -l"))
	pipeline.AddStage(createCommandStageWithName("second", "ls -la"))
	spvar := NewSecialVariables(pipeline)
	spvar.Replace("__RESULT[\"first\"]")
}
