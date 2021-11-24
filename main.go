package main

import (
	"log"
	"rss-reader/cmd"
)

func main() {
	if err := cmd.Execute("PORT", "LOGIN", "PASSWORD", "RSS_FETCH_EVERY_N_SECS", "DATABASE_URL"); err {
		log.Fatalf("Unable to start commands: %v", err)
	}
}
