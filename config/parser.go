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
	"container/list"
	"reflect"
	"strings"

	"github.com/recruit-tech/plumber/log"
	"github.com/recruit-tech/plumber/pipelines"
	"github.com/recruit-tech/plumber/stages"
)

func getStageTypeModuleName(stageType string) string {
	return strings.ToLower(stageType)
}

func Parse(configData *map[interface{}]interface{}) *pipelines.Pipeline {
	pipeline := pipelines.NewPipeline()
	pipelineData, ok := (*configData)["pipeline"].([]interface{})
	if ok == false {
		return pipeline
	}

	stageList := convertYamlMapToStages(pipelineData)
	for stageItem := stageList.Front(); stageItem != nil; stageItem = stageItem.Next() {
		pipeline.AddStage(stageItem.Value.(stages.Stage))
	}
	return pipeline
}

func convertYamlMapToStages(yamlStageList []interface{}) *list.List {
	stages := list.New()
	for _, stageDetail := range yamlStageList {
		stage := mapStage(stageDetail.(map[interface{}]interface{}))
		stages.PushBack(stage)
	}
	return stages
}

func mapStage(stageMap map[interface{}]interface{}) stages.Stage {
	log.Debugf("%v", stageMap["run_after"])
	stageType := stageMap["stage_type"].(string)
	stage := stages.InitStage(stageType)
	newStageValue := reflect.ValueOf(stage).Elem()
	newStageType := reflect.TypeOf(stage).Elem()

	if stageName := stageMap["stage_name"]; stageName != nil {
		stage.SetStageName(stageMap["stage_name"].(string))
	}

	for i := 0; i < newStageType.NumField(); i++ {
		tagName := newStageType.Field(i).Tag.Get("config")
		for stageOptKey, stageOptVal := range stageMap {
			if tagName == stageOptKey {
				fieldVal := newStageValue.Field(i)
				if fieldVal.Type() == reflect.ValueOf("string").Type() {
					fieldVal.SetString(stageOptVal.(string))
				}
			}
		}
	}

	if runAfters := stageMap["run_after"]; runAfters != nil {
		for _, runAfter := range runAfters.([]interface{}) {
			childStage := mapStage(runAfter.(map[interface{}]interface{}))
			stage.AddChildStage(childStage)
		}
	}
	return stage
}
