package config

import (
	"reflect"
	"strings"

	"github.com/takahi-i/plumber/pipelines"
	"github.com/takahi-i/plumber/stages"
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
