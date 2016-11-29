package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/walter-cd/walter/lib/pipeline"
)

func main() {
	const defaultConfigFile = "pipeline.yml"

	var (
		configFile string
		version    bool
		build      bool
		deploy     bool
	)

	flag.StringVar(&configFile, "c", defaultConfigFile, "file which define pipeline")
	flag.BoolVar(&version, "v", false, "print version string")
	flag.BoolVar(&build, "build", false, "run build")
	flag.BoolVar(&deploy, "deploy", false, "run deploy")

	flag.Parse()

	if version {
		log.Info(OutputVersion())
		os.Exit(0)
	}

	if !build && !deploy {
		log.Error("specify -build and/or -deploy flags")
		os.Exit(1)
	}

	p, err := pipeline.LoadFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(p.Run(build, deploy))
}
