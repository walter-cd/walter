/* plumber: a deployment pipeline template
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

	"github.com/recruit-tech/plumber/stages"
)

func TestParse(t *testing.T) {
	configData := ReadConfig("../tests/fixtures/pipeline.yml")
	actual := (*Parse(configData)).Stages.Front().Value.(*stages.CommandStage).Command

	expected := "echo \"hello, world\""
	t.Logf("got %v\nwant %v", actual, expected)
	//if expected != actual {
	//t.Errorf("got %v\nwant %v", actual, expected)
	//}
}
