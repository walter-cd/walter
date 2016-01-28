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

//Package messengers provides all functionality for the suported messengers
package messengers

import (
	"fmt"
)

// Messenger is a interface for notifying the result to the messeging
// services such as Slack or HipChat.
type Messenger interface {
	Post(string, ...string) bool
	Suppress(string) bool
}

//BaseMessenger struct
type BaseMessenger struct {
	SuppressFields []string `config:"suppress"`
}

//Post posts the supplied message
func (baseMessenger *BaseMessenger) Post(messege string) bool {
	return true
}

//Suppress idenitifies if the supplied output type is a suppressed field
func (baseMessenger *BaseMessenger) Suppress(outputType string) bool {
	for _, suppress := range baseMessenger.SuppressFields {
		if suppress == outputType {
			return true
		}
	}
	return false
}

// InitMessenger generates a spefified Messenger client objet.
func InitMessenger(mtype string) (Messenger, error) {
	var messenger Messenger
	switch mtype {
	case "hipchat":
		messenger = new(HipChat)
	case "hipchat2":
		messenger = new(HipChat2)
	case "slack":
		messenger = new(Slack)
	case "fake":
		messenger = new(FakeMessenger)
	default:
		err := fmt.Errorf("no messenger type: %s", mtype)
		return nil, err
	}
	return messenger, nil
}
