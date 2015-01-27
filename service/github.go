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
	"container/list"
	"time"
	"io/ioutil"
	"encoding/json"

	"github.com/recruit-tech/walter/log"
	"github.com/google/go-github/github"
	"code.google.com/p/goauth2/oauth"
)

type Update struct {
	Time time.Time `json:"time"`
	Succeeded bool `json:"succeeded"`
	Status string  `json:"status"`
}

type GitHubClient struct {
	Repo string `config:"repo"`
	From string `config:"from"`
	Token string `config:"token"`
	UpdateFile string `config:"update"`
	Update
}

func (self *GitHubClient) GetCommits() (*list.List, error) {
	log.Info("getting commits\n");
	commits := list.New()
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: self.Token},
	}
	client := github.NewClient(t.Client())

	// get a list of pull requests with Pull Request API
	pullreqs, _, err := client.PullRequests.List(self.From, self.Repo,
			&github.PullRequestListOptions{})
	if err != nil {
		log.Errorf("failed to get pull requests");
		return list.New(), err
	}

	for _, pullreq := range pullreqs {
		if *pullreq.State == "Open" && pullreq.UpdatedAt.After(self.Time) {
			commits.PushBack(pullreq)
		}
	}

	// get the latest commit with Commit API
	master_commits, _, _ := client.Repositories.ListCommits(
	self.From, self.Repo, &github.CommitsListOptions{})
	if master_commits[0].Commit.Author.Date.After(self.Time) {
		commits.PushBack(master_commits[0])
	}
	return commits, nil
}

func LoadLastUpdate(fname string) (Update, error) {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return Update{}, err
	}
	log.Infof("loading last update form %s\n", string(file));
	var update Update
	if err:= json.Unmarshal(file, &update); err != nil {
		return Update{}, err
	}
	return update, nil
}

func SaveUpdate(fname string, update Update) bool {
	log.Infof("writing new update form %s\n", string(fname));
	bytes, err:= json.Marshal(update)
	if err != nil {
		log.Errorf("failed to convert update to string...: %s\n", err.Error());
		return false
	}
		if err:= ioutil.WriteFile(fname, bytes, 644); err != nil {
		log.Errorf("failed to write update to file...: %s\n", err.Error());
	}
	return false
}
