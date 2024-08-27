package main

import "fmt"

//nolint:gochecknoglobals
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printVersion() {
	fmt.Printf("version: %s\ncommit: %s\ndate: %s\n", version, commit, date)
}
