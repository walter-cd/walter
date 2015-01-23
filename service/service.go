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
)

type Service interface {
	Run(Commit) Result
	GetCommits() list.List
	Load(fname string)
}

type Commit struct {
	Updated bool
	SHA1   string
	Branch string
	Date time.Time
}

type Result struct {
	Message string
	Success bool
	Date time.Time
}

// func InitService(stype string) (Service, error) {
// 	var service Service
// 	switch stype {
// 	case "github":
// 		service = new(GitHub)
// 	default:
// 		err := fmt.Errorf("no service type: %s", stype)
// 		return nil, err
// 	}
// 	return service, nil
// }
