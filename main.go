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
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/recruit-tech/walter/config"
	"github.com/recruit-tech/walter/log"
	"github.com/recruit-tech/walter/version"
	"github.com/recruit-tech/walter/walter"
)

func main() {
	log.Init(&log.GlogRecorder{})

	opts, err := config.LoadOpts(os.Args[1:])

	switch err {
	case nil:
	case flag.ErrHelp:
		os.Exit(0)
	default:
		os.Exit(2)
	}

	if opts.PrintVersion {
		fmt.Printf("Walter version %s\n", version.Version)
		os.Exit(0)
	}

	walter, err := walter.New(opts)
	if err != nil {
		log.Error(err.Error())
		log.Error("failed to create Walter")
		return
	}
	log.Info("running Walter")
	result := walter.Run()
	if result == false {
		log.Error("more than one failures were detected running Walter")
		log.Flush()
		os.Exit(1)
	}
	log.Info("succeded to finish Walter")
	log.Flush()
}
