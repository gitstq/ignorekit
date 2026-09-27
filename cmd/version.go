package cmd

import "fmt"

func init() {
	register(command{name: "version", short: "print version", run: runVersion})
}

func runVersion(args []string) int {
	fmt.Printf("ignorekit version %s\ncommit: %s\nbuild date: %s\n", Version, Commit, BuildDate)
	return exitOK
}
