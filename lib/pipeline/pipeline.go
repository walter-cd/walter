package pipeline

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/walter-cd/walter/lib/task"
)

type Pipeline struct {
	Build  Build
	Deploy Deploy
}

type Build struct {
	Tasks   []task.Task
	Cleanup []task.Task
}

type Deploy struct {
	Tasks   []task.Task
	Cleanup []task.Task
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

func (p Pipeline) Run() {
	p.runBuild()
	p.runDeploy()
}

func (p Pipeline) runBuild() {
	for _, t := range p.Build.Tasks {
		t.Run()
	}

	for _, c := range p.Build.Cleanup {
		c.Run()
	}
}

func (p Pipeline) runDeploy() {
	for _, t := range p.Deploy.Tasks {
		t.Run()
	}

	for _, c := range p.Deploy.Cleanup {
		c.Run()
	}
}
