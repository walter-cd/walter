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
	"regexp"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
	"github.com/recruit-tech/walter/log"
)

//GitHubClient struct
type GitHubClient struct {
	Repo         string `config:"repo"`
	From         string `config:"from"`
	Token        string `config:"token"`
	UpdateFile   string `config:"update"`
	TargetBranch string `config:"branch"`
}

//GetUpdateFilePath returns the update file name
func (githubClient *GitHubClient) GetUpdateFilePath() string {
	if githubClient.UpdateFile != "" {
		return githubClient.UpdateFile
	}
	return DEFAULT_UPDATE_FILE_NAME

}

//RegisterResult registers the supplied result
func (githubClient *GitHubClient) RegisterResult(result Result) error {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: githubClient.Token},
	}
	client := github.NewClient(t.Client())

	log.Info("Submitting result")
	repositories := client.Repositories
	status, _, err := repositories.CreateStatus(
		githubClient.From,
		githubClient.Repo,
		result.SHA,
		&github.RepoStatus{
			State:       github.String(result.State),
			TargetURL:   github.String(""),
			Description: github.String(result.Message),
			Context:     github.String("continuous-integraion/walter"),
		})
	log.Infof("Submit status: %s", status)
	if err != nil {
		log.Errorf("Failed to register result: %s", err)
	}
	return err
}

//GetCommits get a list of all the commits for the current update
func (githubClient *GitHubClient) GetCommits(update Update) (*list.List, error) {
	log.Info("getting commits\n")
	commits := list.New()
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: githubClient.Token},
	}
	client := github.NewClient(t.Client())

	// get a list of pull requests with Pull Request API
	pullreqs, _, err := client.PullRequests.List(
		githubClient.From, githubClient.Repo,
		&github.PullRequestListOptions{})
	if err != nil {
		log.Errorf("Failed to get pull requests")
		return list.New(), err
	}

	re, err := regexp.Compile(githubClient.TargetBranch)
	if err != nil {
		log.Error("Failed to compile branch pattern...")
		return list.New(), err
	}

	log.Infof("Size of pull reqests: %d", len(pullreqs))
	for _, pullreq := range pullreqs {
		log.Infof("Branch name is \"%s\"", *pullreq.Head.Ref)

		if githubClient.TargetBranch != "" {
			matched := re.Match([]byte(*pullreq.Head.Ref))
			if matched != true {
				log.Infof("Not add a branch, \"%s\" since this branch name is not match the filtering pattern", *pullreq.Head.Ref)
				continue
			}
		}

		if *pullreq.State == "open" && pullreq.UpdatedAt.After(update.Time) {
			log.Infof("Adding pullrequest %d", *pullreq.Number)
			commits.PushBack(pullreq)
		}
	}

	// get the latest commit with Commit API if the commit is newer than last update
	masterCommits, _, _ := client.Repositories.ListCommits(
		githubClient.From, githubClient.Repo, &github.CommitsListOptions{})
	if masterCommits[0].Commit.Author.Date.After(update.Time) {
		commits.PushBack(masterCommits[0])
	}
	return commits, nil
}
