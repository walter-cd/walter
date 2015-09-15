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
package engine

import (
	"regexp"

	"github.com/recruit-tech/walter/pipelines"
)

// SpecialVariables is a set of variables contains all the results and
// outputs of previous stages.
type SpecialVariables struct {
	pipeline *pipelines.Pipeline
	re       *regexp.Regexp
}

// Replace replaces all environment variables in a line
func (specialVariables *SpecialVariables) Replace(line string) (string, error) {
	// Search list of start and stop positions of special variables

	for result := (*specialVariables.re).FindStringSubmatchIndex(line); result != nil; result = (*specialVariables.re).FindStringSubmatchIndex(line) {
		// Extract spacial variable of specified stages
		value, err := specialVariables.extractStageValue(line, result)
		if err != nil {
			return "", err
		}
		line = line[0:result[0]] + value + line[result[1]:]
	}
	return line, nil
}

func (specialVariables *SpecialVariables) extractStageValue(line string, pos []int) (string, error) {
	outType := line[pos[2]:pos[3]]
	stageName := line[pos[4]:pos[5]]
	return specialVariables.pipeline.GetStageResult(stageName, outType)
}

func NewSecialVariables(pipeline *pipelines.Pipeline) *SpecialVariables {
	regexStr := "(__RESULT|__OUT|__ERR)\\[\"([a-zA-Z_]+)\"\\]"
	pattern, _ := regexp.Compile(regexStr)
	return &SpecialVariables{
		pipeline: pipeline,
		re:       pattern,
	}
}
