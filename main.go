package main

import (
	"log"
	"rss-reader/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Unable to start commands: %v", err)
	}
}
