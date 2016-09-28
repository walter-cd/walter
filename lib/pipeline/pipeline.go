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
	b.Tasks.Run()
	b.Cleanup.Run()
}

func (d *Deploy) Run() {
	d.Tasks.Run()
	d.Cleanup.Run()
}
