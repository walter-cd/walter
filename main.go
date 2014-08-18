package main

import (
	"fmt"
	"os"

	"github.com/takahi-i/plumber/config"
	"github.com/takahi-i/plumber/plumber"
)

func main() {
	fmt.Printf("Running plumber\n")
	opts := config.LoadOpts(os.Args[1:])
	var plumber = plumber.New(opts)
	plumber.Run()
}
