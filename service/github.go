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
	"time"
	"io/ioutil"
	"encoding/json"

	"github.com/recruit-tech/walter/log"
	//"github.com/google/go-github/github"
	//"code.google.com/p/goauth2/oauth"
)

type GitHub struct {
	Repository string `config:"repository"`
	From  string `config:"from"`
	Token string `config:"token"`
}

type Update struct {
	Time time.Time `json:"time"`
	Succeeded bool `json:"succeeded"`
	Status string  `json:"status"`
}

type Client struct {
	github GitHub
	update Update
}

func LoadLastUpdate(fname string) (Update, error) {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return Update{}, err
	}
	log.Infof("Loading last update form %s\n", string(file));
	var update Update
	if err:= json.Unmarshal(file, &update); err != nil {
		return Update{}, err
	}
	return update, nil
}

func SaveUpdate(fname string, update Update) bool {
	log.Infof("Writing new update form %s\n", string(fname));
	bytes, err:= json.Marshal(update)
	if err != nil {
		log.Errorf("Failed to convert update to string...: %s\n", err.Error());
		return false
	}
	if err:= ioutil.WriteFile(fname, bytes, 644); err != nil {
		log.Errorf("Failed to write update to file...: %s\n", err.Error());
	}
	return false
}
