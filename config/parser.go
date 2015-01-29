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
	"container/list"
	"fmt"
	"reflect"
	"strings"

	"github.com/recruit-tech/walter/log"
	"github.com/recruit-tech/walter/messengers"
	"github.com/recruit-tech/walter/pipelines"
	"github.com/recruit-tech/walter/services"
	"github.com/recruit-tech/walter/stages"
)

func getStageTypeModuleName(stageType string) string {
	return strings.ToLower(stageType)
}

// TODO: need refactoring
func Parse(configData *map[interface{}]interface{}) (*pipelines.Pipeline, error) {
	// parse service block
	serviceOps, ok := (*configData)["service"].(map[interface{}]interface{})
	var repoService services.Service
	var err error
	if ok == true {
		log.Info("Found \"service\" block")
		repoService, err = mapService(serviceOps)
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("Not found \"service\" block")
		repoService, err = services.InitService("local")
		if err != nil {
			return nil, err
		}
	}

	// parse messenger block
	messengerOps, ok := (*configData)["messenger"].(map[interface{}]interface{})
	var messenger messengers.Messenger
	if ok == true {
		log.Info("Found messenger block")
		messenger, err = mapMessenger(messengerOps)
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("Not found messenger block")
		messenger, err = messengers.InitMessenger("fake")
		if err != nil {
			return nil, err
		}
	}

	// parse pipeline block
	var pipeline *pipelines.Pipeline = &pipelines.Pipeline{
		Reporter: messenger, RepoService: repoService,
	}

	pipelineData, ok := (*configData)["pipeline"].([]interface{})
	if ok == false {
		return nil, fmt.Errorf("No pipeline block in the input file")
	}
	stageList, err := convertYamlMapToStages(pipelineData)
	if err != nil {
		return nil, err
	}
	for stageItem := stageList.Front(); stageItem != nil; stageItem = stageItem.Next() {
		pipeline.AddStage(stageItem.Value.(stages.Stage))
	}
	return pipeline, nil
}

func mapMessenger(messengerMap map[interface{}]interface{}) (messengers.Messenger, error) {
	messengerType := messengerMap["type"].(string)
	log.Info("type of reporter is " + messengerType)
	messenger, err := messengers.InitMessenger(messengerType)
	if err != nil {
		return nil, err
	}
	newMessengerValue := reflect.ValueOf(messenger).Elem()
	newMessengerType := reflect.TypeOf(messenger).Elem()
	for i := 0; i < newMessengerType.NumField(); i++ {
		tagName := newMessengerType.Field(i).Tag.Get("config")
		for messengerOptKey, messengerOptVal := range messengerMap {
			if tagName == messengerOptKey {
				fieldVal := newMessengerValue.Field(i)
				if fieldVal.Type() == reflect.ValueOf("string").Type() {
					fieldVal.SetString(messengerOptVal.(string))
				}
			}
		}
	}

	return messenger, nil
}

func mapService(serviceMap map[interface{}]interface{}) (services.Service, error) {
	serviceType := serviceMap["type"].(string)
	log.Info("type of service is " + serviceType)
	service, err := services.InitService(serviceType)
	if err != nil {
		return nil, err
	}

	newServiceValue := reflect.ValueOf(service).Elem()
	newServiceType := reflect.TypeOf(service).Elem()
	for i := 0; i < newServiceType.NumField(); i++ {
		tagName := newServiceType.Field(i).Tag.Get("config")
		for serviceOptKey, serviceOptVal := range serviceMap {
			if tagName == serviceOptKey {
				fieldVal := newServiceValue.Field(i)
				if fieldVal.Type() == reflect.ValueOf("string").Type() {
					fieldVal.SetString(serviceOptVal.(string))
				}
			}
		}
	}
	return service, nil
}

func convertYamlMapToStages(yamlStageList []interface{}) (*list.List, error) {
	stages := list.New()
	for _, stageDetail := range yamlStageList {
		stage, err := mapStage(stageDetail.(map[interface{}]interface{}))
		if err != nil {
			return nil, err
		}
		stages.PushBack(stage)
	}
	return stages, nil
}

func mapStage(stageMap map[interface{}]interface{}) (stages.Stage, error) {
	log.Debugf("%v", stageMap["run_after"])

	var stageType string = "command"
	if stageMap["stage_type"] != nil {
		stageType = stageMap["stage_type"].(string)
	}
	stage, err := stages.InitStage(stageType)
	if err != nil {
		return nil, err
	}
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
			childStage, err := mapStage(runAfter.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			stage.AddChildStage(childStage)
		}
	}
	return stage, nil
}
