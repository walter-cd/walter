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
package walter

import (
	"os/exec"
	"reflect"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"github.com/recruit-tech/walter/config"
	"github.com/recruit-tech/walter/engine"
	"github.com/recruit-tech/walter/log"
	"github.com/recruit-tech/walter/services"
	"github.com/recruit-tech/walter/stages"
)

// Walter object.
type Walter struct {
	Engine *engine.Engine
	Opts   *config.Opts
}

// New creates a Walter instance.
func New(opts *config.Opts) (*Walter, error) {
	log.Infof("Pipeline file path: \"%s\"", opts.PipelineFilePath)
	configData := config.ReadConfig(opts.PipelineFilePath)
	resources, err := config.Parse(configData)
	if err != nil {
		return nil, err
	}
	monitorCh := make(chan stages.Mediator)
	engine := &engine.Engine{
		Resources: resources,
		Opts:      opts,
		MonitorCh: &monitorCh,
	}
	return &Walter{
		Opts:   opts,
		Engine: engine,
	}, nil
}

// Run executes registered Pipeline.
func (e *Walter) Run() bool {
	repoServiceValue := reflect.ValueOf(e.Engine.Resources.RepoService)
	if e.Engine.Opts.Mode == "local" ||
		repoServiceValue.Type().String() == "*services.LocalClient" {
		log.Info("Starting Walter in local mode")
		result := e.Engine.RunOnce()
		return result.IsSucceeded()
	} else {
		log.Info("Starting Walter in repository service mode")
		return e.runService()
	}
}

func (e *Walter) runService() bool {
	// load .walter-update
	log.Infof("Loading update file... \"%s\"", e.Engine.Resources.RepoService.GetUpdateFilePath())
	update, err := services.LoadLastUpdate(e.Engine.Resources.RepoService.GetUpdateFilePath())
	log.Infof("Succeeded loading update file")

	log.Info("Updating status...")
	update.Status = "inprogress"
	result := services.SaveLastUpdate(e.Engine.Resources.RepoService.GetUpdateFilePath(), update)
	if result == false {
		log.Error("Failed to save status update")
		return false
	}
	log.Info("Succeeded updating status")

	// get latest commit and pull requests
	log.Info("downloading commits and pull requests...")
	commits, err := e.Engine.Resources.RepoService.GetCommits(update)
	if err != nil {
		log.Errorf("Failed getting commits: %s", err)
		return false
	}

	log.Info("Succeeded getting commits")
	log.Info("Size of commits: " + strconv.Itoa(commits.Len()))
	has_failed_process := false
	for commit := commits.Front(); commit != nil; commit = commit.Next() {
		commitType := reflect.TypeOf(commit.Value)
		if commitType.Name() == "RepositoryCommit" {
			log.Info("Found new repository commit")
			trunkCommit := commit.Value.(github.RepositoryCommit)
			if result := e.processTrunkCommit(trunkCommit); result == false {
				has_failed_process = true
			}
		} else if commitType.Name() == "PullRequest" {
			log.Info("Found new pull request commit")
			pullreq := commit.Value.(github.PullRequest)
			if result := e.processPullRequest(pullreq); result == false {
				has_failed_process = true
			}
		} else {
			log.Errorf("Nothing commit type: %s", commitType)
			has_failed_process = true
		}
	}

	// save .walter-update
	log.Info("Saving update file...")
	update.Status = "finished"
	update.Time = time.Now()
	result = services.SaveLastUpdate(e.Engine.Resources.RepoService.GetUpdateFilePath(), update)
	if result == false {
		log.Error("Failed to save update")
		return false
	}
	return !has_failed_process
}

func (e *Walter) processTrunkCommit(commit github.RepositoryCommit) bool {
	log.Infof("Checkout master branch")
	_, err := exec.Command("git", "checkout", "master", "-f").Output()
	if err != nil {
		log.Errorf("Failed to checkout master branch: %s", err)
		return false
	}
	log.Infof("Downloading new commit from master")
	_, err = exec.Command("git", "pull", "origin", "master").Output()
	if err != nil {
		log.Errorf("Failed to download new commit from master: %s", err)
		return false
	}
	log.Infof("Running the latest commit in master")
	w, err := New(e.Opts)
	if err != nil {
		log.Errorf("Failed to create Walter object...: %s", err)
		log.Error("Skip execution...")
		return false
	}
	result := w.Engine.RunOnce()

	// register the result to hosting service
	if result.IsSucceeded() {
		log.Info("Succeeded.")
		e.Engine.Resources.RepoService.RegisterResult(
			services.Result{
				State:   "success",
				Message: "Succeeded running pipeline...",
				SHA:     *commit.SHA})
		return true
	} else {
		log.Error("Error reported...")
		e.Engine.Resources.RepoService.RegisterResult(
			services.Result{
				State:   "failure",
				Message: "Failed running pipleline ...",
				SHA:     *commit.SHA})
		return false
	}
}

func (e *Walter) processPullRequest(pullrequest github.PullRequest) bool {
	// checkout pullrequest
	num := *pullrequest.Number
	_, err := exec.Command("git", "fetch", "origin", "refs/pull/"+strconv.Itoa(num)+"/head:pr_"+strconv.Itoa(num)).Output()

	defer exec.Command("git", "checkout", "master", "-f").Output() // TODO: make trunk branch configurable
	defer log.Info("returning master branch...")

	if err != nil {
		log.Errorf("Failed to fetch pull request: %s", err)
		return false
	}

	_, err = exec.Command("git", "checkout", "pr_"+strconv.Itoa(num)).Output()
	if err != nil {
		log.Errorf("Failed to checkout pullrequest branch (\"%s\") : %s", "pr_"+strconv.Itoa(num), err)
		log.Error("Skip execution...")
		return false
	}

	// run pipeline
	log.Info("Running pipeline...")
	w, err := New(e.Opts)
	if err != nil {
		log.Errorf("Failed to create Walter object...: %s", err)
		log.Error("Skip execution...")
		return false
	}

	result := w.Engine.RunOnce()

	// register the result to hosting service
	if result.IsSucceeded() {
		log.Info("succeeded.")
		e.Engine.Resources.RepoService.RegisterResult(
			services.Result{
				State:   "success",
				Message: "Succeeded running pipeline...",
				SHA:     *pullrequest.Head.SHA})
		return true
	} else {
		log.Error("Error reported...")
		e.Engine.Resources.RepoService.RegisterResult(
			services.Result{
				State:   "failure",
				Message: "Failed running pipleline ...",
				SHA:     *pullrequest.Head.SHA})
		return false
	}
}
