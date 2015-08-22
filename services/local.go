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

import "container/list"

//LocalClient struct
type LocalClient struct{}

//RegisterResult resgisters the supplied result
func (localClient *LocalClient) RegisterResult(result Result) error {
	return nil
}

//GetCommits gets the commits for the current update
func (localClient *LocalClient) GetCommits(update Update) (*list.List, error) {
	return list.New(), nil
}

//GetUpdateFilePath returns the update file path
func (localClient *LocalClient) GetUpdateFilePath() string {
	return ""
}
