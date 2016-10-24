package pipeline

import (
	"io/ioutil"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"

	"github.com/go-yaml/yaml"
	"github.com/walter-cd/walter/lib/task"
)

type Pipeline struct {
	Build  Build
	Deploy Deploy
}

type Build struct {
	Tasks   task.Tasks
	Cleanup task.Tasks
}

type Deploy struct {
	Tasks   task.Tasks
	Cleanup task.Tasks
}

func Load(b []byte) (Pipeline, error) {
	p := Pipeline{}
	err := yaml.Unmarshal(b, &p)
	return p, err
}

func LoadFromFile(file string) (Pipeline, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return Pipeline{}, err
	}
	return Load(data)
}

func (p *Pipeline) Run() {
	p.Build.Run()
	//p.Deploy.Run()
}

func (b *Build) Run() {
	log.Info("Start build phase")
	buildCtx, buildCancel := context.WithCancel(context.Background())
	defer buildCancel()
	b.Tasks.Run(buildCtx, buildCancel)

	log.Info("Start cleanup phase of build")
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()
	b.Cleanup.Run(cleanupCtx, cleanupCancel)
}

func (d *Deploy) Run() {
	log.Info("Start deploy phase")
	deployCtx, deployCancel := context.WithCancel(context.Background())
	defer deployCancel()
	d.Tasks.Run(deployCtx, deployCancel)

	log.Info("Start cleanup phase of deploy")
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()
	d.Cleanup.Run(cleanupCtx, cleanupCancel)
}
