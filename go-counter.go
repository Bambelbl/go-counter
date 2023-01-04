package main

import (
	"github.com/Bambelbl/go-counter/pkg/processor"
	"log"
)

const k = 5

func main() {
	p, err := processor.NewProcessor(k)
	if err != nil {
		log.Fatalf("Error in create new processor with k = %d: %s", k, err.Error())
	}
	err = p.Process()
	if err != nil {
		log.Fatalf("Error in process sources: %s", err.Error())
	}
}
