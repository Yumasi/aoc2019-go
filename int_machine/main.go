package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments provided")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open file: %s", err)
	}

	if mem, err := newMemory(file); err != nil {
		log.Fatalf("could not init memory from file %s: %s", os.Args[1], err)
	} else if err = mem.Exec(); err != nil {
		log.Fatalf("an error occured during execution: %s\nInput: %v", err, mem.memory)
	}
}
