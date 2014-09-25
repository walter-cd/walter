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
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	configData := *ReadConfig("../tests/fixtures/pipeline.yml")
	actual := configData["pipeline"].([]interface{})[0].(map[interface{}]interface{})["command"]
	assert.Equal(t, "echo \"hello, world\"", actual)
}

func TestReadConfigBytes(t *testing.T) {
	configStr :=
		`pipeline:
    - stage_name: command_stage_1
      stage_type: command
      command: echo "hello, world"
`
	configBytes := []byte(configStr)
	configData := *ReadConfigBytes(configBytes)
	actual := configData["pipeline"].([]interface{})[0].(map[interface{}]interface{})["command"]

	assert.Equal(t, "echo \"hello, world\"", actual)
}

func TestReadConfigWithChildren(t *testing.T) {
	configStr :=
		`pipeline:
    - stage_name: command_stage_1
      stage_type: command
      command: echo "hello, world"
      run_after:
          -  stage_name: command_stage_2_group_1
             stage_type: command
             command: echo "hello, world, command_stage_2_group_1"
    - stage_name: command_stage_3
      stage_type: command
      command: echo "hello, world"1
`
	configBytes := []byte(configStr)
	configData := *ReadConfigBytes(configBytes)
	pipelineConf := configData["pipeline"].([]interface{})[0].(map[interface{}]interface{})
	actual := pipelineConf["run_after"].([]interface{})[0].(map[interface{}]interface{})["command"]
	assert.Equal(t, "echo \"hello, world, command_stage_2_group_1\"", actual)
}

func TestReadPipelineWithoutStageConfig(t *testing.T) {
	configStr := "pipeline:"
	configBytes := []byte(configStr)
	configData := *ReadConfigBytes(configBytes)
	actual, _ := configData["pipeline"]
	assert.Nil(t, actual)
}

func TestReadVoidConfig(t *testing.T) {
	configStr := ""
	configBytes := []byte(configStr)
	configData := *ReadConfigBytes(configBytes)
	actual := len(configData)
	assert.Equal(t, 0, actual)
}
