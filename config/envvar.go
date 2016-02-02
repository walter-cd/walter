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

	"github.com/walter-cd/walter/log"
)

// EnvVariables is a set of environment variables contains all the variables
// defined when the walter command is executed.
type EnvVariables struct {
	variables  *map[string]string
	envPattern *regexp.Regexp
	spPattern  *regexp.Regexp
}

// NewEnvVariables creates one EnvVariable object.
func NewEnvVariables() *EnvVariables {
	envmap := loadEnvMap()
	envPattern, _ := regexp.Compile("[$]([a-zA-Z_]+)")
	spPattern, _ := regexp.Compile("(__RESULT|__OUT|__ERR|__COMBINED)\\[\"([a-zA-Z_0-9 ]+)\"\\]")

	return &EnvVariables{
		variables:  &envmap,
		envPattern: envPattern,
		spPattern:  spPattern,
	}
}

// Get returns the value of envionment variable.
func (envVariables *EnvVariables) Get(vname string) (string, bool) {
	replaced := envVariables.replaceSpecialVariable(vname)
	val, ok := (*envVariables.variables)[replaced]
	return val, ok
}

// Add appends the value to specified envionment variable.
func (envVariables *EnvVariables) Add(key string, value string) {
	(*envVariables.variables)[key] = value
}

// ExportSpecialVariable appends the value of special variable as a envionment variable.
func (envVariables *EnvVariables) ExportSpecialVariable(key string, value string) {
	replaced := envVariables.replaceSpecialVariable(key)
	(*envVariables.variables)[replaced] = value
	os.Setenv(replaced, value) //NOTE: export environment variable
}

// Replace replaces all environment variables in a line
func (envVariables *EnvVariables) Replace(line string) string {
	converted := envVariables.ReplaceSpecialVariableToEnvVariable(line)
	ret := (*envVariables.envPattern).ReplaceAllStringFunc(converted, envVariables.regexReplace)
	return ret
}

// ReplaceSpecialVariableToEnvVariable only special variables in the given string to environment variable representation
func (envVariables *EnvVariables) ReplaceSpecialVariableToEnvVariable(line string) string {
	for result := (*envVariables.spPattern).FindStringSubmatchIndex(line); result != nil; result = (*envVariables.spPattern).FindStringSubmatchIndex(line) {
		outType := line[result[2]:result[3]]
		stageName := line[result[4]:result[5]]
		stageName = strings.Replace(stageName, " ", "_", -1)
		line = line[0:result[0]] + "$" + outType + "__" + stageName + "__" + line[result[1]:]
	}
	return line
}

func (envVariables *EnvVariables) replaceSpecialVariable(key string) string {
	pos := (*envVariables.spPattern).FindStringSubmatchIndex(key)
	if pos == nil {
		return key
	}
	outType := key[pos[2]:pos[3]]
	stageName := key[pos[4]:pos[5]]
	stageName = strings.Replace(stageName, " ", "_", -1)
	return outType + "__" + stageName + "__"
}

func (envVariables *EnvVariables) regexReplace(input string) string {
	matched := (*envVariables.envPattern).FindStringSubmatch(input)
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
