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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvAccess(t *testing.T) {
	envs := NewEnvVariables()
	path, ok := envs.Get("GOPATH")
	assert.True(t, ok)
	assert.True(t, len(path) > 0)
}

func TestEnvAccessNoExist(t *testing.T) {
	envs := NewEnvVariables()
	_, ok := envs.Get("NO_SUCH_A_ENV_VARIABLE")
	assert.False(t, ok)
}

func TestReplaceLineWithEnvVariable(t *testing.T) {
	envs := NewEnvVariables()
	envs.Add("SLACK_CHANNEL", "foobar")
	result := envs.Replace("path: $SLACK_CHANNEL")
	assert.Equal(t, "path: foobar", result)
}

func TestReplaceMultipleItemsLineWithEnvVariable(t *testing.T) {
	envs := NewEnvVariables()
	envs.Add("PATH", "/usr/:/usr/local")
	envs.Add("LOCAL", "en")
	result := envs.Replace("$PATH is set for $LOCAL")
	assert.Equal(t, "/usr/:/usr/local is set for en", result)
}

func TestReplaceWithoutWhiteSpace(t *testing.T) {
	envs := NewEnvVariables()
	envs.Add("PATH", "/usr/:/usr/local")
	result := envs.Replace("$PATH:/opt")
	assert.Equal(t, "/usr/:/usr/local:/opt", result)
}

func TestVoidInput(t *testing.T) {
	envs := NewEnvVariables()
	result := envs.Replace("")
	assert.Equal(t, "", result)
}
