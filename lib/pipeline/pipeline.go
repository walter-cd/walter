package pipeline

import (
	"io/ioutil"

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
	p.Deploy.Run()
}

func (b *Build) Run() {
	log.Info("Start build phase")
	b.Tasks.Run()
	log.Info("Start cleanup phase of build")
	b.Cleanup.Run()
}

func (d *Deploy) Run() {
	log.Info("Start deploy phase")
	d.Tasks.Run()
	log.Info("Start cleanup phase of build")
	d.Cleanup.Run()
}
