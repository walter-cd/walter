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
	"github.com/walter-cd/walter/messengers"
	"github.com/walter-cd/walter/services"
	"github.com/walter-cd/walter/stages"
)

func TestParseFromFile(t *testing.T) {
	configData, err := ReadConfig("../tests/fixtures/pipeline.yml")
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	actual := resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello, world\"", actual)
	assert.Nil(t, err)
}

func TestParseJustHeading(t *testing.T) {
	configData, err := ReadConfigBytes([]byte("pipeline:"))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	pipeline, err := parser.Parse()
	assert.Nil(t, pipeline)
	assert.NotNil(t, err)
}

func TestParseVoid(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(""))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	pipeline, err := parser.Parse()
	assert.Nil(t, pipeline)
	assert.NotNil(t, err)
}

func TestParseConfWithChildren(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      type: command
      command: echo "hello, world"
      run_after:
          -  name: command_stage_2_group_1
             type: command
             command: echo "hello, world, command_stage_2_group_1"
          -  name: command_stage_3_group_1
             type: command
             command: echo "hello, world, command_stage_3_group_1"`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	assert.Equal(t, 1, result.Pipeline.Size())
	assert.Nil(t, err)

	childStages := result.Pipeline.Stages.Front().Value.(stages.Stage).GetChildStages()
	assert.Equal(t, 2, childStages.Len())
}

func TestParseConfWithParallel(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: parallel stages
      parallel:
          -  name: parallel command 1
             type: command
             command: echo "hello, world, parallel command 1"
          -  name: parallel command 2
             type: command
             command: echo "hello, world, parallel command 2"`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	assert.Equal(t, 1, result.Pipeline.Size())
	assert.Nil(t, err)

	childStages := result.Pipeline.Stages.Front().Value.(stages.Stage).GetChildStages()
	assert.Equal(t, 2, childStages.Len())
}

func TestParseConfDefaultStageTypeIsCommand(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      command: echo "hello, world"
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	actual := result.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello, world\"", actual)
	assert.Nil(t, err)
}

func TestParseConfWithDirectory(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      type: command
      command: ls -l
      directory: /usr/local
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	actual := result.Pipeline.Stages.Front().Value.(*stages.CommandStage).Directory
	assert.Nil(t, err)
	assert.Equal(t, "/usr/local", actual)
}

func TestParseConfWithShellScriptStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      type: shell
      file: ../stages/test_sample.sh
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	actual := result.Pipeline.Stages.Front().Value.(*stages.ShellScriptStage).File
	assert.Equal(t, "../stages/test_sample.sh", actual)
	assert.Nil(t, err)
}

func TestParseConfWithMessengerBlock(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
    messenger:
           type: hipchat
           room_id: foobar
           token: xxxx
           from: yyyy
    pipeline:
        - name: command_stage_1
          type: shell
          file: ../stages/test_sample.sh
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	messenger, ok := result.Reporter.(*messengers.HipChat)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "foobar", messenger.RoomID)
	assert.Equal(t, "xxxx", messenger.Token)
	assert.Equal(t, "yyyy", messenger.From)
}

func TestParseConfWithMessengerBlockWithSupress(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
    messenger:
           type: hipchat
           room_id: foobar
           token: xxxx
           from: yyyy
           suppress:
              - stderr
              - stdout

    pipeline:
        - name: command_stage_1
          type: shell
          file: ../stages/test_sample.sh
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	messenger, ok := result.Reporter.(*messengers.HipChat)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "foobar", messenger.RoomID)
	assert.Equal(t, "xxxx", messenger.Token)
	assert.Equal(t, "yyyy", messenger.From)
	assert.Equal(t, 2, len(messenger.SuppressFields))
}

func TestParseConfWithInvalidStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      type: xxxxx
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestParseConfWithInvalidChildStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      type: command
      command: echo "hello, world"
      run_after:
          -  name: command_stage_2_group_1
             type: xxxxx
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestParseConfWithInvalidParallelStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: parallel stages
      parallel:
          -  name: parallel command 1
             type: xxxxx
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestParseConfWithServiceBlock(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
    service:
        type: github
        token: xxxx
        repo: walter
        from: yyyy
        update: .walter-update
    pipeline:
        - name: command_stage_1
          type: shell
          file: ../stages/test_sample.sh
    `))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	service, ok := result.RepoService.(*services.GitHubClient)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "xxxx", service.Token)
	assert.Equal(t, "walter", service.Repo)
	assert.Equal(t, "yyyy", service.From)
}

func TestParseConfWithEnvVariable(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      command: echo "hello $USER_NAME"
`))
	assert.Nil(t, err)
	envs := NewEnvVariables()
	envs.Add("USER_NAME", "takahi-i")
	parser := &Parser{ConfigData: configData, EnvVariables: envs}
	result, err := parser.Parse()
	actual := result.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello $USER_NAME\"", actual)
	assert.Nil(t, err)
}

func TestParseSpecialVariableWithSpace(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: stage 1
      command: echo "hello foobar"
    - name: stage 2
      command: echo hello __OUT["stage 1"]
`))
	assert.Nil(t, err)
	envs := NewEnvVariables()
	parser := &Parser{ConfigData: configData, EnvVariables: envs}
	result, err := parser.Parse()
	actual := result.Pipeline.Stages.Back().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo hello $__OUT__stage_1__", actual)
	assert.Nil(t, err)
}

func TestParseConfWithSpecialVariable(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: stage_1
      command: echo "hello foobar"
    - name: stage_2
      command: echo hello __OUT["stage_1"]
`))
	assert.Nil(t, err)
	envs := NewEnvVariables()
	parser := &Parser{ConfigData: configData, EnvVariables: envs}
	result, err := parser.Parse()
	actual := result.Pipeline.Stages.Back().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo hello $__OUT__stage_1__", actual)
	assert.Nil(t, err)
}

func TestParseConfWithEnvVariableInDirectoryAttribute(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      command: echo "hello walter"
      directory: $HOME
`))
	assert.Nil(t, err)
	envs := NewEnvVariables()
	parser := &Parser{ConfigData: configData, EnvVariables: envs}
	_, err2 := parser.Parse() // confirm not to be panic
	assert.Nil(t, err2)
}

func TestParseConfWithNoExistEnvVariable(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - name: command_stage_1
      command: echo "hello $NO_SUCH_A_ENV_VARIABLE"
`))
	assert.Nil(t, err)
	envs := NewEnvVariables()
	envs.Add("USER_NAME", "takahi-i")
	parser := &Parser{ConfigData: configData, EnvVariables: envs}
	result, err := parser.Parse()
	actual := result.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello $NO_SUCH_A_ENV_VARIABLE\"", actual) // NOTE: No env variable name is shown when there is no env variable
	assert.Nil(t, err)
}

func TestParseMessengerConfWithEnvVariable(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
    messenger:
           type: hipchat
           room_id: foobar
           token: $HIPCHAT_TOKEN
           from: yyyy
    pipeline:
        - name: command_stage_1
          type: shell
          file: ../stages/test_sample.sh
`))
	assert.Nil(t, err)
	envs := NewEnvVariables()
	envs.Add("HIPCHAT_TOKEN", "this-token-is-very-secret")
	parser := &Parser{ConfigData: configData, EnvVariables: envs}
	result, err := parser.Parse()
	messenger, ok := result.Reporter.(*messengers.HipChat)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "foobar", messenger.RoomID)
	assert.Equal(t, "this-token-is-very-secret", messenger.Token)
	assert.Equal(t, "yyyy", messenger.From)
}

func TestParseConfigWithDeprecatedProperties(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`pipeline:
    - stage_name: command_stage_1
      stage_type: command
      command: echo "hello, world"
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	assert.Equal(t, 1, result.Pipeline.Size())
	assert.Nil(t, err)
	actual := result.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello, world\"", actual)
}

func TestParseFromFileWithRequire(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
require:
    - ../tests/fixtures/s2_stages.yml

pipeline:
  - call: s2::foo
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.Equal(t, "echo \"hello foo in s2\"", resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command)
	assert.Equal(t, "s2::foo", resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).GetStageName())

	assert.Nil(t, err)
}

func TestParseFromFileWithRequireStageWithEnv(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
require:
    - ../tests/fixtures/s1_stages.yml

pipeline:
  - call: s1::foo
`))
	assert.Nil(t, err)
	envs := NewEnvVariables()
	envs.Add("VAR1", "Heroku")
	parser := &Parser{ConfigData: configData, EnvVariables: envs}
	resources, err := parser.Parse()
	actual := resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello foo with VAR1=$VAR1 in s1\"", actual)
	assert.Nil(t, err)
}

func TestParseFromFileWithRequireNonExistStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
require:
    - ../tests/fixtures/s1_stages.yml

pipeline:
  - call: s1::foobar
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.Nil(t, resources)
	assert.NotNil(t, err)
}

func TestParseFromFileWithRequiredCleanUpStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
require:
    - ../tests/fixtures/s2_stages.yml

pipeline:
  - call: s2::foo

cleanup:
  - call: s2::bar
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.Equal(t, "echo \"hello bar in s2\"", resources.Cleanup.Stages.Front().Value.(*stages.CommandStage).Command)
	assert.Equal(t, "s2::bar", resources.Cleanup.Stages.Front().Value.(*stages.CommandStage).GetStageName())
	assert.Nil(t, err)
}

func TestParseFromFileWithRequiredParallelStages(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
require:
    - ../tests/fixtures/s2_stages.yml

pipeline:
  - name: parallel stages
    parallel:
        - call: s2::foo
        - call: s2::bar
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.Nil(t, err)
	childStages := resources.Pipeline.Stages.Front().Value.(stages.Stage).GetChildStages()
	assert.Equal(t, 2, childStages.Len())
	childStage1 := childStages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello foo in s2\"", childStage1)
	childStage2 := childStages.Back().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo \"hello bar in s2\"", childStage2)
}

func TestParseFromFileWithAddingFeatureToRequiredStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
require:
    - ../tests/fixtures/s1_stages.yml

pipeline:
  - call: s1::foo
    directory: /
`))
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.NotNil(t, resources)
	assert.Nil(t, err)
}

func TestParseFromFileWithInvalidOverrideRequiredStage(t *testing.T) {
	configData, err := ReadConfigBytes([]byte(`
require:
    - ../tests/fixtures/s1_stages.yml

pipeline:
  - call: s1::foo
    command: echo "hello, world"
`))

	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.Nil(t, resources)
	assert.NotNil(t, err)
	assert.Equal(t, "overriding required stage is forbidden", err.Error())
}
