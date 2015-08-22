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
	"os"
	"regexp"
	"strings"

	"github.com/recruit-tech/walter/log"
)

// EnvVariables is a set of environment variables contains all the variables
// defined when the walter command is executed.
type EnvVariables struct {
	variables *map[string]string
	re        *regexp.Regexp
}

// NewEnvVariables creates one EnvVariable object.
func NewEnvVariables() *EnvVariables {
	envmap := loadEnvMap()
	regexStr := "[$]([a-zA-Z_]+)"
	envPattern, _ := regexp.Compile(regexStr)
	return &EnvVariables{
		variables: &envmap,
		re:        envPattern,
	}
}

// Get returns the value of envionment variable.
func (envVariables *EnvVariables) Get(vname string) (string, bool) {
	val, ok := (*envVariables.variables)[vname]
	return val, ok
}

// Add appends the value to specified envionment variable.
func (envVariables *EnvVariables) Add(key string, value string) {
	(*envVariables.variables)[key] = value
}

// Replace replaces all environment variables in a line
func (envVariables *EnvVariables) Replace(line string) string {
	ret := (*envVariables.re).ReplaceAllStringFunc(line, envVariables.regexReplace)
	return ret
}

func (envVariables *EnvVariables) regexReplace(input string) string {
	matched := (*envVariables.re).FindStringSubmatch(input)
	if len(matched) == 2 {
		if replaced := (*envVariables.variables)[matched[1]]; replaced != "" {
			return replaced
		}
		log.Warnf("NO environment variable: %s", matched[0])
		return ""

	}
	return input
}

func loadEnvMap() map[string]string {
	envs := make(map[string]string)
	for _, envVal := range os.Environ() {
		curEnv := strings.Split(envVal, "=")
		envs[curEnv[0]] = curEnv[1]
	}
	return envs
}
