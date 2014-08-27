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
package stages

import (
	"testing"
)

func TestWIthSimpleCommand(t *testing.T) {
	stage := NewCommandStage()
	stage.AddCommand("ls", "-l")
	expected := true
	actual := stage.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestWithNoexistCommand(t *testing.T) {
	stage := NewCommandStage()
	stage.AddCommand("zzzz", "")
	expected := false
	actual := stage.Run()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestStdoutRsultOfCommand(t *testing.T) {
	stage := NewCommandStage()
	stage.AddCommand("echo", "foobar")
	expected := "foobar\n"
	stage.Run()
	actual := stage.GetStdoutResult()
	if expected != actual {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
