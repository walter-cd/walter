package main

import "fmt"
import "plumber"

func main() {
	fmt.Printf("Running plumber\n");
	pipeline := plumber.NewPipeline()
	pipeline.AddStage(plumber.NewCommandStage())
}
