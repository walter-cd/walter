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
package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitStage(t *testing.T) {
	stage, err := InitStage("command")
	assert.NotNil(t, stage)
	assert.Nil(t, err)
}

func TestInitNonExistStage(t *testing.T) {
	stage, err := InitStage("xxxx")
	assert.Nil(t, stage)
	assert.NotNil(t, err)
}

// TODO: simplify the test case
func TestAddChildStage(t *testing.T) {
	stage := &CommandStage{}
	stage.StageName = "test_command_stage"
	PrepareCh(stage)

	child := &CommandStage{}
	child.StageName = "test_child"
	PrepareCh(child)

	stage.AddCommand("ls -l")
	child.AddCommand("ls -l")

	stage.AddChildStage(child)
	childStages := stage.GetChildStages()
	assert.Equal(t, 1, childStages.Len())
}
