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

	"github.com/recruit-tech/plumber/stages"
)

func TestParseFromFile(t *testing.T) {
	configData := ReadConfig("../tests/fixtures/pipeline.yml")
	actual := (*Parse(configData)).Stages.Front().Value.(*stages.CommandStage).Command

	expected := "echo \"hello, world\""
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestParseJustHeading(t *testing.T) {
	configData := ReadConfigBytes([]byte("pipeline:"))
	actual := Parse(configData)

	expected := 0
	if expected != actual.Size() {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestParseVoid(t *testing.T) {
	configData := ReadConfigBytes([]byte(""))
	actual := Parse(configData).Size()

	expected := 0
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
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
             command: echo "hello, world, command_stage_3_group_1"
`))
	result := Parse(configData)

	expectedPipelineSize := 1
	if expectedPipelineSize != result.Size() {
		t.Errorf("got %v\nwant %v", result.Size(), expectedPipelineSize)
	}

	childStages := result.Stages.Front().Value.(stages.Stage).GetChildStages()

	expectedChildStageSize := 2
	if expectedChildStageSize != childStages.Len() {
		t.Errorf("got %v\nwant %v", childStages.Len(), expectedChildStageSize)
	}

}
