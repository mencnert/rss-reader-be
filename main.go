package main

import (
	"log"
	"rss-reader/cmd"
)

func main() {
	if err := cmd.Execute(
		"PORT",
		"LOGIN",
		"PASSWORD",
		"RSS_FETCH_EVERY_N_SECS",
		"CLEAN_DB_EVERY_N_HOURS",
		"DATABASE_URL",
	); err != nil {
		log.Fatalf("Unable to start commands: %v", err)
	}
}
