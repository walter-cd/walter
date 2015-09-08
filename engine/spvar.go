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
	"container/list"
	"fmt"
	"regexp"

	"github.com/recruit-tech/walter/pipelines"
)

// SpecialVariables is a set of variables contains all the results and
// outputs of previous stages.
type SpecialVariables struct {
	// Pipeline
	pipeline *pipelines.Pipeline
	re       *regexp.Regexp
}

// Replace replaces all environment variables in a line
func (specialVariables *SpecialVariables) Replace(line string) string {
	// Search list of start and stop positions of special variables
	specialVariables.extractPositions(line)
	// Extract stage names
	// Extract spacial variable of specified stages
	return line
}

func (specialVariables *SpecialVariables) extractPositions(line string) *list.List {
	ret := list.New()
	results := (*specialVariables.re).FindAllStringSubmatch(line, -1)

	for _, result := range results {
		fmt.Println("found match")
		fmt.Println(result)
	}
	return ret
}

func NewSecialVariables(pipeline *pipelines.Pipeline) *SpecialVariables {
	regexStr := "(__RESULT|__OUT|__ERR)\\[\"([a-zA-Z_]+)\"\\]"
	pattern, err := regexp.Compile(regexStr)
	if err != nil {
		fmt.Println("Failed to compile regex..")
	}

	return &SpecialVariables{
		pipeline: pipeline,
		re:       pattern,
	}
}
