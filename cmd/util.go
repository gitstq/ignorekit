package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/gitstq/ignorekit/internal/term"
)

func outf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func errf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func bold(s string) string { return term.Bold(s) }

// maxList caps how many items are printed per section in terminal output.
const maxList = 15

// commonFlags builds a flag set shared across commands and returns the
// option holders plus the set.
type commonOpts struct {
	dir     string
	format  string
	write   bool
	untrack bool
	yes     bool
	noColor bool
}

func bindCommon(fs *flag.FlagSet, o *commonOpts) {
	fs.StringVar(&o.dir, "C", ".", "project directory")
	fs.StringVar(&o.dir, "dir", ".", "project directory")
	fs.StringVar(&o.format, "f", "text", "output format: text|json|sarif")
	fs.StringVar(&o.format, "format", "text", "output format: text|json|sarif")
	fs.BoolVar(&o.write, "write", false, "write to .gitignore")
	fs.BoolVar(&o.untrack, "untrack", false, "also git rm --cached tracked ignored files")
	fs.BoolVar(&o.yes, "y", false, "proceed without prompting")
	fs.BoolVar(&o.yes, "yes", false, "proceed without prompting")
	fs.BoolVar(&o.noColor, "no-color", false, "disable colored output")
}

func applyColor(o *commonOpts) {
	if o.noColor {
		term.Enabled = false
	}
}

func parse(fs *flag.FlagSet, args []string) error {
	return fs.Parse(args[1:])
}
