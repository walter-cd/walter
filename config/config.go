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

// Package config defines the configration parameters,
// and the parser to load configuration file.
package config

import (
	"flag"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

var (
	fs = flag.NewFlagSet("walter", flag.ExitOnError)
)

// Opts contains a set of configuration options.
type Opts struct {
	PipelineFilePath string
	StopOnAnyFailure bool
	PrintVersion     bool
	Mode             string
}

// LoadOpts defines the prameters of the walter command.
func LoadOpts(arguments []string) (*Opts, error) {
	var pipelineFilePath string
	var stopOnAnyFailure bool
	var printVersion bool
	var threshold string
	var log_dir string
	var mode string

	fs.StringVar(&pipelineFilePath, "c", "./pipeline.yml", "pipeline.yml file")
	fs.BoolVar(&stopOnAnyFailure, "f", false, "Skip execution of subsequent stage after failing to exec the upstream stage.")
	fs.BoolVar(&printVersion, "v", false, "Print the version information and exit.")
	fs.StringVar(&threshold, "threshold", "INFO", "Log events at or above this severity are logged.")
	fs.StringVar(&log_dir, "log_dir", "", "Log files will be written to this directory.")
	fs.StringVar(&mode, "mode", "local", "Execution mode (local or service).")

	if err := fs.Parse(arguments); err != nil {
		return nil, err
	}

	flag.CommandLine.Lookup("stderrthreshold").Value.Set(threshold)

	if log_dir != "" {
		flag.CommandLine.Lookup("log_dir").Value.Set(log_dir)
	}

	return &Opts{
		PipelineFilePath: pipelineFilePath,
		StopOnAnyFailure: stopOnAnyFailure,
		PrintVersion:     printVersion,
		Mode:             mode,
	}, nil
}

func ReadConfig(configFilePath string) (*map[interface{}]interface{}, error) {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	return ReadConfigBytes(data)
}

func ReadConfigBytes(configSetting []byte) (*map[interface{}]interface{}, error) {
	configData := make(map[interface{}]interface{})
	err := yaml.Unmarshal(configSetting, &configData)
	if err != nil {
		return nil, err
	}
	return &configData, nil
}
