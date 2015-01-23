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
package service

import (
  "os"
  "testing"
  "time"
  "io/ioutil"

  "github.com/stretchr/testify/assert"
)

func TestLoadUpdate(t *testing.T) {
  update, err := LoadLastUpdate("update_sample.json")  // NOTE the time format need to be RFC3339
  assert.Nil(t, err)
  assert.NotNil(t, update)
  assert.Equal(t, true, update.Succeeded)
  assert.Equal(t, "finished", update.Status)
  expectedTime, _ := time.Parse(time.RFC3339, "2015-01-21T05:05:42Z")
  assert.Equal(t, expectedTime, update.Time)
}

func TestSaveUpdate(t *testing.T) {
  // save update
  time, _ := time.Parse(time.RFC3339, "2015-01-21T05:05:42Z")
  updateSample := Update{time, true, "finished"}
  tempFile, _ := ioutil.TempFile("", "update")
  defer os.Remove(tempFile.Name())
  defer tempFile.Close()
  SaveUpdate(tempFile.Name(), updateSample)

  // load update
  loadedUpdate, err := LoadLastUpdate(tempFile.Name())
  assert.Nil(t, err)
  assert.NotNil(t, loadedUpdate)
  assert.Equal(t, true, loadedUpdate.Succeeded)
  assert.Equal(t, "finished", loadedUpdate.Status)
}
