package config

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/takahi-i/plumber/third_party/gopkg.in/yaml.v1"
)

type Opts struct {
	PipelineFilePath string
}

func LoadOpts(arguments []string) *Opts {
	var pipelineFilePath string
	flag.StringVar(&pipelineFilePath, "c", "./pipeline.yml", "pipeline.yml file")
	flag.Parse()

	return &Opts{
		PipelineFilePath: pipelineFilePath,
	}
}

func ReadConfig(configFilePath string) *map[interface{}]interface{} {
	configData := make(map[interface{}]interface{})
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("error :%v \n", err)
	}

	err = yaml.Unmarshal([]byte(data), &configData)
	if err != nil {
		fmt.Printf("error :%v \n", err)
	}
	return &configData
}
