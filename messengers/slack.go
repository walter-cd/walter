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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/recruit-tech/walter/log"
)

// Slack is a client which reports the pipeline results to the Slack chennel.
type Slack struct {
	BaseMessenger `config:"suppress"`
	Channel       string `config:"channel" json:"channel"`
	UserName      string `config:"username" json:"username"`
	IconEmoji     string `config:"icon" json:"icon_emoji,omitempty"`
	IconURL       string `config:"icon_url" json:"icon_url,omitempty"`
	IncomingURL   string `config:"url" json:"-"` // not map to json
}

//FakeSlack To avoid the infinite recursion
// (see http://stackoverflow.com/questions/23045884/can-i-use-marshaljson-to-add-arbitrary-fields-to-a-json-encoding-in-golang)
type FakeSlack Slack

//Post posts a message to slack
func (slack *Slack) Post(message string) bool {
	if slack.Channel[0] != '#' {
		log.Infof("Add # to channel name: %s", slack.Channel)
		slack.Channel = "#" + slack.Channel
	}

	var color string

	if strings.Contains(message, "[RESULT] Failed") {
		color = "danger"
	} else if strings.Contains(message, "[RESULT] Skipped") {
		color = "warning"
	} else if strings.Contains(message, "[RESULT] Succeeded") {
		color = "good"
	}

	attachment := map[string]string{
		"text":  message,
		"color": color,
	}

	attachments := []map[string]string{attachment}

	params, _ := json.Marshal(struct {
		FakeSlack
		Attachments []map[string]string `json:"attachments"`
	}{
		FakeSlack:   FakeSlack(*slack),
		Attachments: attachments,
	})

	resp, err := http.PostForm(
		slack.IncomingURL,
		url.Values{"payload": {string(params)}},
	)
	defer resp.Body.Close()

	if err != nil {
		log.Errorf("Failed post message to Slack...: %s", message)
		return false
	}

	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Infof("Slack post result...: %s", body)
		return true
	}
	log.Errorf("Failed to read result from Slack...")
	return false
}
