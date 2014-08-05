package main

import (
	"fmt"

	"github.com/takahi-i/plumber/plumber"
)

func main() {
	fmt.Printf("Running plumber\n")
	var plumber = plumber.New()
	plumber.Run()
}
