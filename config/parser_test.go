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

	"github.com/recruit-tech/walter/messengers"
	"github.com/recruit-tech/walter/services"
	"github.com/recruit-tech/walter/stages"
	"github.com/stretchr/testify/assert"
)

func TestParseFromFile(t *testing.T) {
	configData := ReadConfig("../tests/fixtures/pipeline.yml")
	pipeline, err := Parse(configData)
	actual := pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello, world\"", actual)
	assert.Nil(t, err)
}

func TestParseJustHeading(t *testing.T) {
	configData := ReadConfigBytes([]byte("pipeline:"))
	pipeline, err := Parse(configData)
	assert.Nil(t, pipeline)
	assert.NotNil(t, err)
}

func TestParseVoid(t *testing.T) {
	configData := ReadConfigBytes([]byte(""))
	pipeline, err := Parse(configData)
	assert.Nil(t, pipeline)
	assert.NotNil(t, err)
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
	result, err := Parse(configData)
	assert.Equal(t, 1, result.Size())
	assert.Nil(t, err)

	childStages := result.Stages.Front().Value.(stages.Stage).GetChildStages()
	assert.Equal(t, 2, childStages.Len())
}

func TestParseConfDefaultStageTypeIsCommand(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      command: echo "hello, world"
`))
	result, err := Parse(configData)
	actual := result.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello, world\"", actual)
	assert.Nil(t, err)
}

func TestParseConfWithDirectory(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      stage_type: command
      command: ls -l
      directory: /usr/local
`))
	result, err := Parse(configData)
	actual := result.Stages.Front().Value.(*stages.CommandStage).Directory
	assert.Nil(t, err)
	assert.Equal(t, "/usr/local", actual)
}

func TestParseConfWithShellScriptStage(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      stage_type: shell
      file: ../stages/test_sample.sh
`))
	result, err := Parse(configData)
	actual := result.Stages.Front().Value.(*stages.ShellScriptStage).File
	assert.Equal(t, "../stages/test_sample.sh", actual)
	assert.Nil(t, err)
}

func TestParseConfWithMessengerBlock(t *testing.T) {
	configData := ReadConfigBytes([]byte(`
    messenger:
           type: hipchat
           room_id: foobar
           token: xxxx
           from: yyyy
    pipeline:
        - stage_name: command_stage_1
          stage_type: shell
          file: ../stages/test_sample.sh
`))
	result, err := Parse(configData)
	messenger, ok := result.Reporter.(*messengers.HipChat)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "foobar", messenger.RoomId)
	assert.Equal(t, "xxxx", messenger.Token)
	assert.Equal(t, "yyyy", messenger.From)
}

func TestParseConfWithInvalidStage(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      stage_type: xxxxx
`))
	result, err := Parse(configData)
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestParseConfWithInvalidChildStage(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      stage_type: command
      command: echo "hello, world"
      run_after:
          -  stage_name: command_stage_2_group_1
             stage_type: xxxxx
`))
	result, err := Parse(configData)
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestParseConfWithServiceBlock(t *testing.T) {
	configData := ReadConfigBytes([]byte(`
    service:
        type: github
        token: xxxx
        repo: walter
        from: yyyy
        update: .walter-update
    pipeline:
        - stage_name: command_stage_1
          stage_type: shell
          file: ../stages/test_sample.sh
    `))
	result, err := Parse(configData)
	service, ok := result.RepoService.(*services.GitHubClient)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "xxxx", service.Token)
	assert.Equal(t, "walter", service.Repo)
	assert.Equal(t, "yyyy", service.From)
}

func TestParseConfWithEnvVariable(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      command: echo "hello $USER_NAME"
`))

	envs := NewEnvVariables()
	envs.Add("USER_NAME", "takahi-i")
	result, err := ParseWithSpecifiedEnvs(configData, envs)
	actual := result.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello takahi-i\"", actual)
	assert.Nil(t, err)
}

func TestParseConfWithNoExistEnvVariable(t *testing.T) {
	configData := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      command: echo "hello $NO_SUCH_A_ENV_VARIABLE"
`))

	envs := NewEnvVariables()
	envs.Add("USER_NAME", "takahi-i")
	result, err := ParseWithSpecifiedEnvs(configData, envs)
	actual := result.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello \"", actual) // NOTE: No env variable name is shown when there is no env variable
	assert.Nil(t, err)
}

func TestParseMessengerConfWithEnvVariable(t *testing.T) {
	configData := ReadConfigBytes([]byte(`
    messenger:
           type: hipchat
           room_id: foobar
           token: $HIPCHAT_TOKEN
           from: yyyy
    pipeline:
        - stage_name: command_stage_1
          stage_type: shell
          file: ../stages/test_sample.sh
`))
	envs := NewEnvVariables()
	envs.Add("HIPCHAT_TOKEN", "this-token-is-very-secret")
	result, err := ParseWithSpecifiedEnvs(configData, envs)
	messenger, ok := result.Reporter.(*messengers.HipChat)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "foobar", messenger.RoomId)
	assert.Equal(t, "this-token-is-very-secret", messenger.Token)
	assert.Equal(t, "yyyy", messenger.From)
}
