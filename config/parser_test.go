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
package config

import (
	"testing"

	"github.com/recruit-tech/walter/stages"
	"github.com/stretchr/testify/assert"
)

func TestParseFromFile(t *testing.T) {
	configData := ReadConfig("../tests/fixtures/pipeline.yml")
	actual := (*Parse(configData)).Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello, world\"", actual)
}

func TestParseJustHeading(t *testing.T) {
	configData := ReadConfigBytes([]byte("pipeline:"))
	actual := Parse(configData).Size()
	assert.Equal(t, 0, actual)
}

func TestParseVoid(t *testing.T) {
	configData := ReadConfigBytes([]byte(""))
	actual := Parse(configData).Size()
	assert.Equal(t, 0, actual)
}

func TestParseConfWithChildren(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      stage_type: command
      command: echo "hello, world"
      run_after:
          -  stage_name: command_stage_2_group_1
             stage_type: command
             command: echo "hello, world, command_stage_2_group_1"
          -  stage_name: command_stage_3_group_1
             stage_type: command
             command: echo "hello, world, command_stage_3_group_1"`))
	result := Parse(configData)
	assert.Equal(t, 1, result.Size())

	childStages := result.Stages.Front().Value.(stages.Stage).GetChildStages()
	assert.Equal(t, 2, childStages.Len())
}

func TestParseConfWithDirectory(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      stage_type: command
      command: ls -l
      directory: /usr/local
`))
	result := Parse(configData)
	actual := result.Stages.Front().Value.(*stages.CommandStage).Directory
	assert.Equal(t, "/usr/local", actual)
}
