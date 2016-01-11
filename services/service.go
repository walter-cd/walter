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

//Package services provides the functionality for all supported services (GitHub)
package services

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/walter-cd/walter/log"
)

//Service interface
type Service interface {
	RegisterResult(Result) error
	GetCommits(update Update) (*list.List, error)
	GetUpdateFilePath() string
}

//Update struct
type Update struct {
	Time      time.Time `json:"time"`
	Succeeded bool      `json:"succeeded"`
	Status    string    `json:"status"`
}

//Result struct
type Result struct {
	State   string
	SHA     string
	Message string
}

// DEFAULT_UPDATE_FILE_NAME is the default file name of status of Walter service.
// The file name can be overridden with the configuration file.
const DefaultUpdateFileName string = ".walter"

//LoadLastUpdate loads the last update
func LoadLastUpdate(fname string) (Update, error) {
	file, err := ioutil.ReadFile(fname)
	log.Infof("opening file: \"%s\"...", fname)
	if err != nil {
		log.Warnf("error occured opening file: \"%s\" ...", fname)
		log.Warnf(err.Error())
		log.Warnf("continue the process with the new settings")
		update := Update{
			Time:      time.Date(1970, time.November, 10, 15, 0, 0, 0, time.Local),
			Succeeded: true,
			Status:    "inprogress"}
		return update, nil
	}

	log.Infof("loading last update form \"%s\"\n", string(file))
	var update Update
	if err := json.Unmarshal(file, &update); err != nil {
		log.Warnf("failed to load \"%s\" ...", fname)
		log.Warnf(err.Error())
		log.Warnf("continue the process with the new settings")
		update := Update{
			Time:      time.Now(),
			Succeeded: true,
			Status:    "inprogress"}
		return update, nil
	}

	if update.Status == "inprogress" {
		return Update{}, errors.New("update is currently run in another process")
	}

	log.Info("setting update status into \"inprogress\"...")
	return update, nil
}

//SaveLastUpdate saves the supplied update to the filename
func SaveLastUpdate(fname string, update Update) bool {
	log.Infof("writing down new update: \"%s\"\n", string(fname))
	bytes, err := json.Marshal(update)
	if err != nil {
		log.Errorf("failed to convert update to string...: %s\n", err)
		return false
	}

	if exist := fileExists(fname); exist == true {
		log.Infof("file exist: \"%s\"", fname)
		if err := os.Remove(fname); err != nil {
			log.Errorf("failed to remove \"%s\" with error: \"%s\"", fname, err)
			return false
		}
		log.Infof("succeeded to remove: \"%s\"", fname)
	}

	if err := ioutil.WriteFile(fname, bytes, 0644); err != nil {
		log.Errorf("failed to write update to file...: \"%s\"\n", err)
		return false
	}
	log.Infof("succeeded to write update file: \"%s\"", fname)
	return true
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//InitService initializes the service
func InitService(stype string) (Service, error) {
	var service Service
	switch stype {
	case "github":
		log.Info("GitHub client was created")
		service = new(GitHubClient)
	case "local":
		log.Info("local client was created")
		service = new(LocalClient)
	default:
		err := fmt.Errorf("no messenger type: %s", stype)
		return nil, err
	}
	return service, nil
}
