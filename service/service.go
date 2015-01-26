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
	"container/list"

	"github.com/recruit-tech/walter/log"
)

type Service interface {
	//Run(list.List)
	GetCommits() *list.List
	LoadLastUpdate(fname string) bool
	SaveLastUpdate(fname string) bool
}

type Result struct {
	Message string
	Success bool
	Date time.Time
}

func InitService(stype string) (Service, error) {
	var service Service
	switch stype {
	case "github":
		service = new(GitHubClient)
	default:
		err := log.Errorf("no service type: %s", stype)
		return nil, err
	}
	return service, nil
}
