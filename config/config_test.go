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

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	configData, err := ReadConfig("../tests/fixtures/pipeline.yml")
	actual := (*configData)["pipeline"].([]interface{})[0].(map[interface{}]interface{})["command"]
	assert.Nil(t, err)
	assert.Equal(t, "echo \"hello, world\"", actual)
}

func TestReadConfigBytes(t *testing.T) {
	configStr :=
		`pipeline:
    - name: command_stage_1
      type: command
      command: echo "hello, world"
`
	configBytes := []byte(configStr)
	configData, err := ReadConfigBytes(configBytes)
	actual := (*configData)["pipeline"].([]interface{})[0].(map[interface{}]interface{})["command"]
	assert.Nil(t, err)
	assert.Equal(t, "echo \"hello, world\"", actual)
}

func TestReadConfigWithDirectory(t *testing.T) {
	configStr :=
		`pipeline:
    - name: command_stage_1
      type: command
      command: echo "hello, world"
      directory: /user/local/bin
`
	configBytes := []byte(configStr)
	configData, err := ReadConfigBytes(configBytes)
	actual := (*configData)["pipeline"].([]interface{})[0].(map[interface{}]interface{})["directory"]
	assert.Nil(t, err)
	assert.Equal(t, "/user/local/bin", actual)
}

func TestReadConfigWithChildren(t *testing.T) {
	configStr :=
		`pipeline:
    - name: command_stage_1
      type: command
      command: echo "hello, world"
      run_after:
          -  name: command_stage_2_group_1
             type: command
             command: echo "hello, world, command_stage_2_group_1"
    - name: command_stage_3
      type: command
      command: echo "hello, world"1
`
	configBytes := []byte(configStr)
	configData, err := ReadConfigBytes(configBytes)
	pipelineConf := (*configData)["pipeline"].([]interface{})[0].(map[interface{}]interface{})
	actual := pipelineConf["run_after"].([]interface{})[0].(map[interface{}]interface{})["command"]
	assert.Equal(t, "echo \"hello, world, command_stage_2_group_1\"", actual)
	assert.Nil(t, err)
}

func TestReadPipelineWithoutStageConfig(t *testing.T) {
	configStr := "pipeline:"
	configBytes := []byte(configStr)
	configData, err := ReadConfigBytes(configBytes)
	actual, _ := (*configData)["pipeline"]
	assert.Nil(t, actual)
	assert.Nil(t, err)
}

func TestReadVoidConfig(t *testing.T) {
	configStr := ""
	configBytes := []byte(configStr)
	configData, err := ReadConfigBytes(configBytes)
	actual := len(*configData)
	assert.Equal(t, 0, actual)
	assert.Nil(t, err)
}
