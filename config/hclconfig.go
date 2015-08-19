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

//Package config includes this file which deals with the option of using
//hashicorp's configuration language (HCL) as an alternative to YAML.
//In addition, HCL provides for both it's own language and JSON files alike,
//making the possibility of both human readable and machine creation of pipelines.
//See https://github.com/hashicorp/hcl for more details
package config

import (
	"github.com/hashicorp/hcl"
	//"github.com/recruit-tech/walter/log"
	"io/ioutil"
)

//HCL2YMLConverter manages the conversion from HCL (or JSON) to YAML
type HCL2YMLConverter struct {
	hclFile       []byte
	hclStructure  map[string]interface{}
	yamlStructure map[interface{}]interface{}
}

//ReadHCLConfig reads the supplied HCL configuration file
func (converter *HCL2YMLConverter) ReadHCLConfig(configFilePath string) (*map[interface{}]interface{}, error) {

	//read the supplied file
	var err error

	converter.hclFile, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	convertedYaml, err := converter.ConvertHCLConfig()
	if err != nil {
		return nil, err
	}

	return convertedYaml, nil

}

//ConvertHCLConfig converts the HCL configuration file supplied
func (converter *HCL2YMLConverter) ConvertHCLConfig() (*map[interface{}]interface{}, error) {

	//initialise the HCL and what will be the resulting YAML equivilent
	converter.hclStructure = make(map[string]interface{})
	converter.yamlStructure = make(map[interface{}]interface{})

	//Decode the HCL file
	err := hcl.Decode(&converter.hclStructure, string(converter.hclFile))
	if err != nil {
		return nil, err
	}

	//Parse the resulting HCL Map
	converter.parseHCL()

	//return the resulting YAML Map
	return &converter.yamlStructure, nil

}

func (converter *HCL2YMLConverter) parseHCL() {

	//is there a pipeline?
	pipeline, pipelineExists := converter.hclStructure["pipeline"].([]map[string]interface{})
	if pipelineExists {
		converter.yamlStructure["pipeline"] = converter.addStages(pipeline, false)
	}

	//is there a cleanup?
	cleanup, cleanupExists := converter.hclStructure["cleanup"].([]map[string]interface{})
	if cleanupExists {
		converter.yamlStructure["cleanup"] = converter.addStages(cleanup, false)
	}

	//Is there a global?
	global, globalExists := converter.hclStructure["global"].([]map[string]interface{})
	if globalExists {
		converter.yamlStructure["global"] = map[interface{}]interface{}{}
		for _, globalsValue := range global {
			converter.yamlStructure["global"] = converter.addParameters(globalsValue)
		}
	}

	//Is there a service?
	service, serviceExists := converter.hclStructure["service"].([]map[string]interface{})
	if serviceExists {
		converter.yamlStructure["service"] = map[interface{}]interface{}{}
		for _, serviceValue := range service {
			converter.yamlStructure["service"] = converter.addParameters(serviceValue)
		}
	}

	//Is there a messenger?
	messenger, messengerExists := converter.hclStructure["messenger"].([]map[string]interface{})
	if messengerExists {
		converter.yamlStructure["messenger"] = map[interface{}]interface{}{}
		for _, messengerValue := range messenger {
			converter.yamlStructure["messenger"] = converter.addParameters(messengerValue)
		}
	}

	//Is there a require?
	requiredFiles, requireExists := converter.hclStructure["require"].([]interface{})
	if requireExists {
		converter.yamlStructure["require"] = requiredFiles
	}

	//is there a namespace
	namespace, namespaceExists := converter.hclStructure["namespace"].(string)
	if namespaceExists {
		converter.yamlStructure["namespace"] = namespace
	}

	//is there a stages
	stages, stagesExists := converter.hclStructure["stages"].([]map[string]interface{})
	if stagesExists {
		converter.yamlStructure["stages"] = converter.addStages(stages, true)
	}

}

func (converter *HCL2YMLConverter) addParameters(parameters map[string]interface{}) map[interface{}]interface{} {
	var params = map[interface{}]interface{}{}
	for paramKey, paramValue := range parameters {
		params[paramKey] = paramValue
	}
	return params
}

func (converter *HCL2YMLConverter) addStages(stages []map[string]interface{}, isDef bool) []interface{} {
	var convertedStages = []interface{}{}
	for _, value := range stages {
		//Based on the YAML (which we followed for now) the definition of stages is defined by a Def
		var searchFor = "stage"
		if isDef {
			searchFor = "def"
		}
		for _, stageValue := range value[searchFor].([]map[string]interface{}) {
			var stage = map[interface{}]interface{}{}
			for stageParamKey, stageParamValue := range stageValue {
				if stageParamKey != "parallel" {
					stage[stageParamKey] = stageParamValue
				} else {
					//it's a parallel stage so need to add the sub stages recursively
					//Notice we do NOT suport the deprecated 'run_after' syntax
					stage["parallel"] = converter.addStages(stageParamValue.([]map[string]interface{}), isDef)
				}
			}
			//append the converted stage to the yaml pipeline structure
			if isDef {
				//it's a definition of a stage, so wrap it up accordingly
				var def = map[interface{}]interface{}{}
				def["def"] = stage
				convertedStages = append(convertedStages, def)
			} else {
				convertedStages = append(convertedStages, stage)
			}
		}
	}
	return convertedStages
}
