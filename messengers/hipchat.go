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
	"github.com/andybons/hipchat"
	"github.com/recruit-tech/walter/log"
)

// HipChat is a client which reports the pipeline results to the HipChat server.
// The client uses V1 of the HipChat API.
type HipChat struct {
	BaseMessenger `config:"suppress"`
	RoomID        string `config:"room_id"`
	Token         string `config:"token"`
	From          string `config:"from"`
}

// Post posts a hipchat message
// TODO: make hipchat api endpoint configurable for on-premises servers
func (hipChat *HipChat) Post(message string) bool {
	client := hipchat.Client{AuthToken: hipChat.Token}
	req := hipchat.MessageRequest{
		RoomId:        hipChat.RoomID,
		From:          hipChat.From,
		Message:       message,
		Color:         hipchat.ColorPurple,
		MessageFormat: hipchat.FormatText,
		Notify:        true,
	}
	if err := client.PostMessage(req); err != nil {
		log.Errorf("Failed post message...: %s", message)
		return false
	}
	return true
}
