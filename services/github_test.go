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
package services

import (
	"testing"
	"time"
	"io/ioutil"
	"os"
	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func createUpdate(timeStr string, success bool, status string) *[]byte {
	updateTime, _ := time.Parse(time.RFC3339, timeStr)
	sample := Update{updateTime, success, status}
	bytes, _ := json.Marshal(sample)
	return &bytes
}

func TestLoadUpdate(t *testing.T) {
	// save update
	tempFile, _ := ioutil.TempFile("", "update")
	defer os.Remove(tempFile.Name())
	bytes := createUpdate("2015-01-21T05:05:42Z", true, "finished")
	ioutil.WriteFile(tempFile.Name(), *bytes, 0644)

	//	load update
	update, err := LoadLastUpdate(tempFile.Name())
	assert.Nil(t, err)
	assert.NotNil(t, update)
	assert.Equal(t, true, update.Succeeded)
	assert.Equal(t, "finished", update.Status)
	expectedTime, _ := time.Parse(time.RFC3339, "2015-01-21T05:05:42Z")
	assert.Equal(t, expectedTime, update.Time)
}

func TestLoadUpdateInProgress(t *testing.T) {
	// save update
	tempFile, _ := ioutil.TempFile("", "update")
	defer os.Remove(tempFile.Name())
	bytes := createUpdate("2015-01-21T05:05:42Z", true, "inprogress")
	ioutil.WriteFile(tempFile.Name(), *bytes, 0644)

	//	load update
	_, err := LoadLastUpdate(tempFile.Name())
	assert.NotNil(t, err)
}

func TestLoadUpdateNotExistPath(t *testing.T) {
	input := "nothing_such_file.json"
	update, err := LoadLastUpdate(input)
	defer os.Remove(input)

	assert.Nil(t, err) // NOTE: continue the process even when there is no specified file
	assert.NotNil(t, update)
	assert.Equal(t, true, update.Succeeded)
	assert.Equal(t, "inprogress", update.Status)
}

func TestSaveLastUpdate(t *testing.T) {
	// save update
	time, _ := time.Parse(time.RFC3339, "2015-01-21T05:05:42Z")
	updateSample := Update{time, true, "finished"}
	tempFile, _ := ioutil.TempFile("", "update")
	defer os.Remove(tempFile.Name())
	SaveLastUpdate(tempFile.Name(), updateSample)

	// load update
	loadedUpdate, err := LoadLastUpdate(tempFile.Name())
	assert.Nil(t, err)
	assert.NotNil(t, loadedUpdate)
	assert.Equal(t, true, loadedUpdate.Succeeded)
	assert.Equal(t, "finished", loadedUpdate.Status)
}
