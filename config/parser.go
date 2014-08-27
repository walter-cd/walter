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
	"reflect"
	"strings"

	"github.com/recruit-tech/plumber/pipelines"
	"github.com/recruit-tech/plumber/stages"
)

func getStageTypeModuleName(stageType string) string {
	return strings.ToLower(stageType)
}

func Parse(configData *map[interface{}]interface{}) *pipelines.Pipeline {
	pipeline := pipelines.NewPipeline()
	pipelineData := (*configData)["pipeline"].(map[interface{}]interface{})

	for _, stageDetail := range pipelineData {
		stageType := stageDetail.(map[interface{}]interface{})["stage_type"].(string)
		stageStruct := stages.InitStage(stageType)
		newStageValue := reflect.ValueOf(stageStruct).Elem()
		newStageType := reflect.TypeOf(stageStruct).Elem()

		for i := 0; i < newStageType.NumField(); i++ {
			tagName := newStageType.Field(i).Tag.Get("config")
			for stageOptKey, stageOptVal := range stageDetail.(map[interface{}]interface{}) {
				if tagName == stageOptKey {
					fieldVal := newStageValue.Field(i)
					if fieldVal.Type() == reflect.ValueOf("string").Type() {
						fieldVal.SetString(stageOptVal.(string))
					}
				}
			}
		}
		pipeline.AddStage(stageStruct)
	}

	return pipeline
}
