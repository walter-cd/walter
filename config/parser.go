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

type Parser struct {
	ConfigData   *map[interface{}]interface{}
	EnvVariables *EnvVariables
}

// Parse reads the specified configuration and create the pipeline.Resource.
func (self *Parser) Parse() (*pipelines.Resources, error) {
	// parse require block
	requires, ok := (*self.ConfigData)["require"].([]interface{})
	if ok == true {
		log.Info("found \"require\" block")
		self.mapRequire(requires)
	} else {
		log.Info("not found \"require\" block")
	}

	// parse service block
	serviceOps, ok := (*self.ConfigData)["service"].(map[interface{}]interface{})
	var repoService services.Service
	var err error
	if ok == true {
		log.Info("found \"service\" block")
		repoService, err = self.mapService(serviceOps)
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("not found \"service\" block")
		repoService, err = services.InitService("local")
		if err != nil {
			return nil, err
		}
	}

	// parse messenger block
	messengerOps, ok := (*self.ConfigData)["messenger"].(map[interface{}]interface{})
	var messenger messengers.Messenger
	if ok == true {
		log.Info("found messenger block")
		messenger, err = self.mapMessenger(messengerOps)
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("not found messenger block")
		messenger, err = messengers.InitMessenger("fake")
		if err != nil {
			return nil, err
		}
	}

	// parse cleanup block
	var cleanup *pipelines.Pipeline = &pipelines.Pipeline{}
	cleanupData, ok := (*self.ConfigData)["cleanup"].([]interface{})
	if ok == true {
		log.Info("found cleanup block")
		cleanupList, err := self.convertYamlMapToStages(cleanupData)
		if err != nil {
			return nil, err
		}
		for stageItem := cleanupList.Front(); stageItem != nil; stageItem = stageItem.Next() {
			cleanup.AddStage(stageItem.Value.(stages.Stage))
		}
	} else {
		log.Info("not found cleanup block in the input file")
	}

	// parse pipeline block
	var pipeline *pipelines.Pipeline = &pipelines.Pipeline{}

	pipelineData, ok := (*self.ConfigData)["pipeline"].([]interface{})
	if ok == false {
		return nil, fmt.Errorf("no pipeline block in the input file")
	}
	stageList, err := self.convertYamlMapToStages(pipelineData)
	if err != nil {
		return nil, err
	}
	for stageItem := stageList.Front(); stageItem != nil; stageItem = stageItem.Next() {
		pipeline.AddStage(stageItem.Value.(stages.Stage))
	}
	var resources = &pipelines.Resources{Pipeline: pipeline, Cleanup: cleanup, Reporter: messenger, RepoService: repoService}

	return resources, nil
}

func (self *Parser) mapRequire(requireList []interface{}) (map[string]interface{}, error) {
	requires := make(map[string]interface{})
	for _, requireFile := range requireList {
		log.Info("register require file: " + requireFile.(string))
		ReadConfig(requireFile.(string))
	}
	return requires, nil
}

func (self *Parser) mapMessenger(messengerMap map[interface{}]interface{}) (messengers.Messenger, error) {
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
					fieldVal.SetString(self.EnvVariables.Replace(messengerOptVal.(string)))
				}
			}
		}
	}

	return messenger, nil
}

func (self *Parser) mapService(serviceMap map[interface{}]interface{}) (services.Service, error) {
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
					fieldVal.SetString(self.EnvVariables.Replace(serviceOptVal.(string)))
				}
			}
		}
	}
	return service, nil
}

func (self *Parser) convertYamlMapToStages(yamlStageList []interface{}) (*list.List, error) {
	stages := list.New()
	for _, stageDetail := range yamlStageList {
		stage, err := self.mapStage(stageDetail.(map[interface{}]interface{}))
		if err != nil {
			return nil, err
		}
		stages.PushBack(stage)
	}
	return stages, nil
}

func (self *Parser) mapStage(stageMap map[interface{}]interface{}) (stages.Stage, error) {
	log.Debugf("%v", stageMap["run_after"])

	var stageType string = "command"
	if stageMap["type"] != nil {
		stageType = stageMap["type"].(string)
	} else if stageMap["stage_type"] != nil {
		log.Warn("found property \"stage_type\"")
		log.Warn("property \"stage_type\" is deprecated. please use \"type\" instead.")
		stageType = stageMap["stage_type"].(string)
	}
	stage, err := stages.InitStage(stageType)
	if err != nil {
		return nil, err
	}
	newStageValue := reflect.ValueOf(stage).Elem()
	newStageType := reflect.TypeOf(stage).Elem()

	if stageName := stageMap["name"]; stageName != nil {
		stage.SetStageName(stageMap["name"].(string))
	} else if stageName := stageMap["stage_name"]; stageName != nil {
		log.Warn("found property \"stage_name\"")
		log.Warn("property \"stage_name\" is deprecated. please use \"stage\" instead.")
		stage.SetStageName(stageMap["stage_name"].(string))
	}

	stageOpts := stages.NewStageOpts()

	if reportingFullOutput := stageMap["report_full_output"]; reportingFullOutput != nil {
		stageOpts.ReportingFullOutput = true
	}

	stage.SetStageOpts(*stageOpts)

	for i := 0; i < newStageType.NumField(); i++ {
		tagName := newStageType.Field(i).Tag.Get("config")
		is_replace := newStageType.Field(i).Tag.Get("is_replace")
		for stageOptKey, stageOptVal := range stageMap {
			if tagName == stageOptKey {
				if stageOptVal == nil {
					log.Warnf("stage option \"%s\" is not specified", stageOptKey)
				} else {
					self.setFieldVal(newStageValue.Field(i), stageOptVal, is_replace)
				}
			}
		}
	}

	parallelStages := stageMap["parallel"]
	if parallelStages == nil {
		if parallelStages = stageMap["run_after"]; parallelStages != nil {
			log.Warn("`run_after' will be obsoleted in near future. Use `parallel' instead.")
		}
	}

	if parallelStages != nil {
		for _, parallelStages := range parallelStages.([]interface{}) {
			childStage, err := self.mapStage(parallelStages.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			stage.AddChildStage(childStage)
		}
	}
	return stage, nil
}

func (self *Parser) setFieldVal(fieldVal reflect.Value, stageOptVal interface{}, is_replace string) {
	if fieldVal.Type() == reflect.ValueOf("string").Type() {
		if is_replace == "true" {
			fieldVal.SetString(self.EnvVariables.Replace(stageOptVal.(string)))
		} else {
			fieldVal.SetString(stageOptVal.(string))
		}
	}
}

func getStageTypeModuleName(stageType string) string {
	return strings.ToLower(stageType)
}
