package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/walter-cd/walter/lib/pipeline"
)

func main() {
	const defaultConfigFile = "pipeline.yml"

	var configFile string
	flag.StringVar(&configFile, "c", defaultConfigFile, "file which define pipeline")
	flag.Parse()

	p, err := pipeline.LoadFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(p.Run())
}
