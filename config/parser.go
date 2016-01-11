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

// Package config defines the configration parameters,
// and the parser to load configuration file.
package config

import (
	"container/list"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/walter-cd/walter/log"
	"github.com/walter-cd/walter/messengers"
	"github.com/walter-cd/walter/pipelines"
	"github.com/walter-cd/walter/services"
	"github.com/walter-cd/walter/stages"
)

// Parser struct used to store config data and environment variables
type Parser struct {
	ConfigData   *map[interface{}]interface{}
	EnvVariables *EnvVariables
}

// Parse reads the specified configuration and create the pipeline.Resource.
func (parser *Parser) Parse() (*pipelines.Resources, error) {
	// parse require block
	requireFiles, ok := (*parser.ConfigData)["require"].([]interface{})
	var required map[string]map[interface{}]interface{}
	var err error
	if ok == true {
		log.Info("found \"require\" block")
		required, err = parser.mapRequires(requireFiles)
		if err != nil {
			log.Error("failed to load requires...")
			return nil, err
		}
		log.Info("number of registered stages: " + strconv.Itoa(len(required)))
	} else {
		log.Info("not found \"require\" block")
	}

	// parse service block
	serviceOps, ok := (*parser.ConfigData)["service"].(map[interface{}]interface{})
	var repoService services.Service
	if ok == true {
		log.Info("found \"service\" block")
		repoService, err = parser.mapService(serviceOps)
		if err != nil {
			log.Error("failed to load service settings...")
			return nil, err
		}
	} else {
		log.Info("not found \"service\" block")
		repoService, err = services.InitService("local")
		if err != nil {
			log.Error("failed to init local mode...")
			return nil, err
		}
	}

	// parse messenger block
	messengerOps, ok := (*parser.ConfigData)["messenger"].(map[interface{}]interface{})
	var messenger messengers.Messenger
	if ok == true {
		log.Info("found messenger block")
		messenger, err = parser.mapMessenger(messengerOps)
		if err != nil {
			log.Error("failed to init messenger...")
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
	var cleanup = &pipelines.Pipeline{}
	cleanupData, ok := (*parser.ConfigData)["cleanup"].([]interface{})
	if ok == true {
		log.Info("found cleanup block")
		cleanupList, err := parser.convertYamlMapToStages(cleanupData, required)
		if err != nil {
			log.Error("failed to create a stage in cleanup...")
			return nil, err
		}
		for stageItem := cleanupList.Front(); stageItem != nil; stageItem = stageItem.Next() {
			cleanup.AddStage(stageItem.Value.(stages.Stage))
		}
	} else {
		log.Info("not found cleanup block in the input file")
	}

	// parse pipeline block
	var pipeline = &pipelines.Pipeline{}

	pipelineData, ok := (*parser.ConfigData)["pipeline"].([]interface{})
	if ok == false {
		return nil, fmt.Errorf("no pipeline block in the input file")
	}
	stageList, err := parser.convertYamlMapToStages(pipelineData, required)
	if err != nil {
		log.Error("failed to create a stage in pipeline...")
		return nil, err
	}
	for stageItem := stageList.Front(); stageItem != nil; stageItem = stageItem.Next() {
		pipeline.AddStage(stageItem.Value.(stages.Stage))
	}
	var resources = &pipelines.Resources{Pipeline: pipeline, Cleanup: cleanup, Reporter: messenger, RepoService: repoService}

	return resources, nil
}

func (parser *Parser) mapRequires(requireList []interface{}) (map[string]map[interface{}]interface{}, error) {
	requires := make(map[string]map[interface{}]interface{})
	for _, requireFile := range requireList {
		replacedFilePath := parser.EnvVariables.Replace(requireFile.(string))
		log.Info("register require file: " + replacedFilePath)
		requireData, err := ReadConfig(replacedFilePath)
		if err != nil {
			return nil, err
		}
		parser.mapRequire(*requireData, &requires)
	}
	return requires, nil
}

func (parser *Parser) mapRequire(requireData map[interface{}]interface{},
	requires *map[string]map[interface{}]interface{}) {
	namespace := requireData["namespace"].(string)
	log.Info("detect namespace: " + namespace)

	stages := requireData["stages"].([]interface{})
	log.Info("number of detected stages: " + strconv.Itoa(len(stages)))

	for _, stageDetail := range stages {
		stageMap := stageDetail.(map[interface{}]interface{})
		for _, values := range stageMap {
			valueMap := values.(map[interface{}]interface{})
			stageKey := namespace + "::" + valueMap["name"].(string)
			log.Info("register stage: " + stageKey)
			(*requires)[stageKey] = valueMap
		}
	}
}

func (parser *Parser) mapMessenger(messengerMap map[interface{}]interface{}) (messengers.Messenger, error) {
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
			if tagName != messengerOptKey {
				continue
			}
			fieldVal := newMessengerValue.Field(i)
			if fieldVal.Type() == reflect.ValueOf("string").Type() {
				fieldVal.SetString(parser.EnvVariables.Replace(messengerOptVal.(string)))
			} else if fieldVal.Type().String() == "messengers.BaseMessenger" {
				elements := messengerOptVal.([]interface{})
				suppressor := fieldVal.Interface().(messengers.BaseMessenger)
				for _, element := range elements {
					suppressor.SuppressFields = append(suppressor.SuppressFields, element.(string))
				}
				fieldVal.Set(reflect.ValueOf(suppressor))
			}
		}
	}
	return messenger, nil
}

func (parser *Parser) mapService(serviceMap map[interface{}]interface{}) (services.Service, error) {
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
			if tagName != serviceOptKey {
				continue
			}
			fieldVal := newServiceValue.Field(i)
			if fieldVal.Type() == reflect.ValueOf("string").Type() {
				fieldVal.SetString(parser.EnvVariables.Replace(serviceOptVal.(string)))
			}
		}
	}
	return service, nil
}

func (parser *Parser) convertYamlMapToStages(yamlStageList []interface{},
	requiredStages map[string]map[interface{}]interface{}) (*list.List, error) {
	stages := list.New()
	for _, stageDetail := range yamlStageList {
		stage, err := parser.mapStage(stageDetail.(map[interface{}]interface{}), requiredStages)
		if err != nil {
			return nil, err
		}
		stages.PushBack(stage)
	}
	return stages, nil
}

func (parser *Parser) mapStage(stageMap map[interface{}]interface{},
	requiredStages map[string]map[interface{}]interface{}) (stages.Stage, error) {
	mergedStageMap, err := parser.extractStage(stageMap, requiredStages)
	if err != nil {
		return nil, err
	}
	stage, err := parser.initStage(mergedStageMap)
	if err != nil {
		return nil, err
	}

	if stageName := mergedStageMap["name"]; stageName != nil {
		stage.SetStageName(mergedStageMap["name"].(string))
	} else if stageName := mergedStageMap["stage_name"]; stageName != nil {
		log.Warn("found property \"stage_name\"")
		log.Warn("property \"stage_name\" is deprecated. please use \"stage\" instead.")
		stage.SetStageName(mergedStageMap["stage_name"].(string))
	}

	stageOpts := stages.NewStageOpts()
	if reportingFullOutput := mergedStageMap["report_full_output"]; reportingFullOutput != nil {
		stageOpts.ReportingFullOutput = true
	}
	stage.SetStageOpts(*stageOpts)

	newStageValue := reflect.ValueOf(stage).Elem()
	newStageType := reflect.TypeOf(stage).Elem()
	for i := 0; i < newStageType.NumField(); i++ {
		tagName := newStageType.Field(i).Tag.Get("config")
		isReplace := newStageType.Field(i).Tag.Get("is_replace")
		for stageOptKey, stageOptVal := range mergedStageMap {
			if tagName != stageOptKey {
				continue
			} else if stageOptVal == nil {
				log.Warnf("stage option \"%s\" is not specified", stageOptKey)
				continue
			}
			parser.setFieldVal(newStageValue.Field(i), stageOptVal, isReplace)
		}
	}

	parallelStages := mergedStageMap["parallel"]
	if parallelStages == nil {
		if parallelStages = mergedStageMap["run_after"]; parallelStages != nil {
			log.Warn("`run_after' will be obsoleted in near future. Use `parallel' instead.")
		}
	}

	if parallelStages != nil {
		for _, parallelStage := range parallelStages.([]interface{}) {
			childStage, err := parser.mapStage(parallelStage.(map[interface{}]interface{}), requiredStages)
			if err != nil {
				return nil, err
			}
			stage.AddChildStage(childStage)
		}
	}
	return stage, nil
}

func (parser *Parser) extractStage(stageMap map[interface{}]interface{},
	requiredStages map[string]map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	if stageMap["call"] == nil {
		return stageMap, nil
	}

	// when "call" is applied
	log.Info("detect call")
	stageName := stageMap["call"].(string)
	calledMap := requiredStages[stageName]
	if calledMap == nil {
		return nil, errors.New(stageName + " is not registerd")
	}
	for fieldName, fieldValue := range calledMap {
		log.Info("fieldName: " + fieldName.(string))
		if _, ok := stageMap[fieldName]; ok {
			return nil, errors.New("overriding required stage is forbidden")
		}
		stageMap[fieldName] = fieldValue
	}
	log.Info("stage name: " + stageName)
	stageMap["name"] = stageName
	return stageMap, nil
}

func (parser *Parser) initStage(stageMap map[interface{}]interface{}) (stages.Stage, error) {
	var stageType = "command"
	if stageMap["type"] != nil {
		stageType = stageMap["type"].(string)
	} else if stageMap["stage_type"] != nil {
		log.Warn("found property \"stage_type\"")
		log.Warn("property \"stage_type\" is deprecated. please use \"type\" instead.")
		stageType = stageMap["stage_type"].(string)
	}
	return stages.InitStage(stageType)
}

func (parser *Parser) setFieldVal(fieldVal reflect.Value, stageOptVal interface{}, isReplace string) {
	if fieldVal.Type() != reflect.ValueOf("string").Type() {
		log.Error("found non string field value type...")
		return
	}
	if isReplace == "true" {
		fieldVal.SetString(parser.EnvVariables.Replace(stageOptVal.(string)))
	} else {
		fieldVal.SetString(parser.EnvVariables.ReplaceSpecialVariableToEnvVariable(stageOptVal.(string)))
	}
}

func getStageTypeModuleName(stageType string) string {
	return strings.ToLower(stageType)
}
