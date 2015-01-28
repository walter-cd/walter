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

	"github.com/recruit-tech/walter/log"
	"github.com/google/go-github/github"
	"code.google.com/p/goauth2/oauth"
)

type GitHubClient struct {
	Repo string `config:"repo"`
	From string `config:"from"`
	Token string `config:"token"`
	UpdateFile string `config:"update"`
}

func (self *GitHubClient) GetUpdateFilePath() string {
	return self.UpdateFile
}

func (self *GitHubClient) GetCommits(update Update) (*list.List, error) {
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

	log.Infof("size of pull reqests: %d", len(pullreqs))
	for _, pullreq := range pullreqs {
		if *pullreq.State == "open" && pullreq.UpdatedAt.After(update.Time) {
			log.Infof("Adding pullrequest %d", *pullreq.Number)
			commits.PushBack(pullreq)
		}
	}

	// get the latest commit with Commit API if the commit is newer than last update
	master_commits, _, _ := client.Repositories.ListCommits(
	self.From, self.Repo, &github.CommitsListOptions{})
	if master_commits[0].Commit.Author.Date.After(update.Time) {
		commits.PushBack(master_commits[0])
	}
	return commits, nil
}
