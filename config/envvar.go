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
	regex_str := "[$]([a-zA-Z_]+)"
	envPattern, _ := regexp.Compile(regex_str)
	return &EnvVariables{
		variables: &envmap,
		re:        envPattern,
	}
}

// Get returns the value of envionment variable.
func (self *EnvVariables) Get(vname string) (string, bool) {
	val, ok := (*self.variables)[vname]
	return val, ok
}

// Add appends the value to specified envionment variable.
func (self *EnvVariables) Add(key string, value string) {
	(*self.variables)[key] = value
}

func (self *EnvVariables) Replace(line string) string {
	ret := (*self.re).ReplaceAllStringFunc(line, self.regexReplace)
	return ret
}

func (self *EnvVariables) regexReplace(input string) string {
	matched := (*self.re).FindStringSubmatch(input)
	if len(matched) == 2 {
		if replaced := (*self.variables)[matched[1]]; replaced != "" {
			return replaced
		} else {
			log.Warnf("NO environment variable: %s", matched[0])
			return ""
		}
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
