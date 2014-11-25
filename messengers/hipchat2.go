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
package messengers

import (
	"github.com/recruit-tech/walter/log"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

// HipChat2 represents a messenger for V2 of the HipChat API
type HipChat2 struct {
	RoomID string `config:"room_id"`
	Token  string `config:"token"`
	From   string `config:"from"`
	client *hipchat.Client
}

// Post sends a new HipChat message using V2 of the API
func (hc *HipChat2) Post(message string) bool {
	if hc.client == nil {
		hc.client = hipchat.NewClient(hc.Token)
	}

	msg := &hipchat.NotificationRequest{
		Color:         "purple",
		Message:       message,
		Notify:        true,
		MessageFormat: "text",
	}

	if _, err := hc.client.Room.Notification(hc.RoomID, msg); err != nil {
		log.Errorf("Failed post message...: %s", msg.Message)
		return false
	}

	return true
}
