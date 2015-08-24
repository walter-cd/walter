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

func TestPipelineFromHCLFile(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	configData, err := hclconverter.ReadHCLConfig("../tests/fixtures/pipeline.hcl")
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	actual := resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo 'hello, world'", actual)
	assert.Nil(t, err)
}

func TestPipelineFromJSONFile(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	configData, err := hclconverter.ReadHCLConfig("../tests/fixtures/pipeline.json")
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	actual := resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command
	assert.Equal(t, "echo 'hello, world'", actual)
	assert.Nil(t, err)
}

func TestParseFromHCLFileWithRequire(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	hclconverter.hclFile = []byte(`
		require = ["../tests/fixtures/s2_stages.hcl"]
		pipeline {
			stage {
				call = "s2::foo"
		  }
		}
`)
	configData, err := hclconverter.ConvertHCLConfig()
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.Equal(t, "echo \"hello foo in s2\"", resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command)
	assert.Equal(t, "s2::foo", resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).GetStageName())

	assert.Nil(t, err)
}

func TestParseFromJSONFileWithRequire(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	hclconverter.hclFile = []byte(`{
		"require":["../tests/fixtures/s2_stages.json"],
		"pipeline":{
			"stage":{
				"call":"s2::foo"
		  }
		}
	}`)
	configData, err := hclconverter.ConvertHCLConfig()
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	resources, err := parser.Parse()
	assert.Equal(t, "echo \"hello foo in s2\"", resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).Command)
	assert.Equal(t, "s2::foo", resources.Pipeline.Stages.Front().Value.(*stages.CommandStage).GetStageName())

	assert.Nil(t, err)
}

func TestParseHCLConfWithMessengerBlock(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	hclconverter.hclFile = []byte(`
	messenger {
	  type = "hipchat"
	  room_id = "foobar"
	  token = "xxxx"
	  from = "yyyy"
	}

  pipeline {
		stage {
			name = "command_stage_1"
			type = "shell"
			file = "../stages/test_sample.sh"
		}
	}
`)
	configData, err := hclconverter.ConvertHCLConfig()
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	messenger, ok := result.Reporter.(*messengers.HipChat)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "foobar", messenger.RoomId)
	assert.Equal(t, "xxxx", messenger.Token)
	assert.Equal(t, "yyyy", messenger.From)
}

func TestParseJSONConfWithMessengerBlock(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	hclconverter.hclFile = []byte(`
	{
		"messenger":{
		  "type":"hipchat",
		  "room_id":"foobar",
		  "token":"xxxx",
		  "from":"yyyy"
		},

	  "pipeline":{
			"stage":{
				"name":"command_stage_1",
				"type":"shell",
				"file":"../stages/test_sample.sh"
			}
		}
	}
`)
	configData, err := hclconverter.ConvertHCLConfig()
	assert.Nil(t, err)
	parser := &Parser{ConfigData: configData, EnvVariables: NewEnvVariables()}
	result, err := parser.Parse()
	messenger, ok := result.Reporter.(*messengers.HipChat)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, "foobar", messenger.RoomId)
	assert.Equal(t, "xxxx", messenger.Token)
	assert.Equal(t, "yyyy", messenger.From)
}

func TestParseHCLConfWithServiceBlock(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	hclconverter.hclFile = []byte(`
	service {
	  type = "github"
		repo = "walter"
	  token = "xxxx"
	  from = "yyyy"
		update = ".walter-update"
	}

  pipeline {
		stage {
			name = "command_stage_1"
			type = "shell"
			file = "../stages/test_sample.sh"
		}
	}
`)
	configData, err := hclconverter.ConvertHCLConfig()
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

func TestParseJSONConfWithServiceBlock(t *testing.T) {
	hclconverter := HCL2YMLConverter{}
	hclconverter.hclFile = []byte(`
	{
		"service":{
		  "type":"github",
		  "repo":"walter",
		  "token":"xxxx",
		  "from":"yyyy",
			"update":".walter-update"
		},
	  "pipeline":{
			"stage":{
				"name":"command_stage_1",
				"type":"shell",
				"file":"../stages/test_sample.sh"
			}
		}
	}
`)
	configData, err := hclconverter.ConvertHCLConfig()
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
