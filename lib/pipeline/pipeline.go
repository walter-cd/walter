package pipeline

import (
	"github.com/go-yaml/yaml"
	"github.com/walter-cd/walter-v2/lib/stage"
)

type Pipeline struct {
	Pipeline []stage.Stage
}

func Load(y string) (Pipeline, error) {
	p := Pipeline{}
	err := yaml.Unmarshal([]byte(y), &p)
	return p, err
}

/*
func LoadFromFile(file string) (Pipeline error) {

}
*/

func (p Pipeline) Run() {

}
