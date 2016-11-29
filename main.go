package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/walter-cd/walter/lib/pipeline"
)

func main() {
	const defaultConfigFile = "pipeline.yml"

	var (
		configFile string
		version    bool
	)

	flag.StringVar(&configFile, "c", defaultConfigFile, "file which define pipeline")
	flag.BoolVar(&version, "v", false, "print version string")

	flag.Parse()

	if version {
		fmt.Println(OutputVersion())
		os.Exit(0)
	}

	p, err := pipeline.LoadFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(p.Run())
}
