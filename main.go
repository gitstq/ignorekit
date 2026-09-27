// Command ignorekit is an offline, zero-dependency .gitignore generator and
// git tracking-hygiene workbench. See the README for full documentation.
package main

import (
	"os"

	"github.com/gitstq/ignorekit/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args))
}
