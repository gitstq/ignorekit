// Package cmd implements the ignorekit command-line interface.
package cmd

import "strings"

// Build-time metadata (overridden via -ldflags).
var (
	Version   = "1.0.0"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Exit codes for CI gating.
const (
	exitOK       = 0 // clean
	exitFindings = 1 // hygiene or secret findings present
	exitError    = 2 // usage / runtime error
)

type command struct {
	name  string
	short string
	run   func(args []string) int
}

var commands = map[string]command{}

func register(c command) { commands[c.name] = c }

// Run dispatches the CLI and returns an exit code.
func Run(argv []string) int {
	in := argv[1:]
	if len(in) == 0 {
		printHelp()
		return exitError
	}

	name, pre, post := splitArgs(in)

	switch name {
	case "":
		printHelp()
		return exitError
	case "help":
		printHelp()
		return exitOK
	case "version":
		full := append([]string{"version"}, append(pre, post...)...)
		return runVersion(full)
	case "-v", "--version":
		return runVersion([]string{"version"})
	case "-h", "--help":
		printHelp()
		return exitOK
	}

	c, ok := commands[name]
	if !ok {
		errf("unknown command %q\n\n", name)
		printHelp()
		return exitError
	}
	full := append([]string{name}, pre...)
	full = append(full, post...)
	return c.run(full)
}

// preTakesValue lists global flags (before the command) that consume a value.
func preTakesValue(f string) bool {
	switch f {
	case "-C", "--dir", "-f", "--format":
		return true
	}
	return false
}

// splitArgs separates the command name (first non-flag token) from global
// flags placed before it and the arguments after it. It understands
// "-f value", "-f=value" and boolean flags.
func splitArgs(in []string) (name string, pre, post []string) {
	i := 0
	for i < len(in) {
		a := in[i]
		if name == "" && !strings.HasPrefix(a, "-") {
			name = a
			i++
			continue
		}
		if name == "" {
			pre = append(pre, a)
			if strings.Contains(a, "=") {
				i++
				continue
			}
			if preTakesValue(a) && i+1 < len(in) {
				pre = append(pre, in[i+1])
				i += 2
				continue
			}
			i++
			continue
		}
		post = append(post, a)
		i++
	}
	return name, pre, post
}

func printHelp() {
	outf(`🦞 %s — modern .gitignore generator & git tracking-hygiene workbench

%s
  init            Detect tools and generate a .gitignore (stdout by default)
  audit           Read-only hygiene report (tracked files that should be ignored)
  apply           Merge suggested rules into .gitignore (safe, idempotent)
  list-templates  List all built-in rule templates
  version         Print version information

%s
  -C, --dir DIR     Project directory (default: current directory)
  -f, --format FMT  Output format: text | json | sarif (audit)
      --write       Write result to .gitignore (init)
      --untrack     Also 'git rm --cached' tracked-but-ignored files (apply)
  -y, --yes         Do not prompt; proceed with changes
      --no-color   Disable colored output
  -h, --help        Show help

%s
  ignorekit init --write
  ignorekit audit --format json
  ignorekit apply --untrack --yes
`, bold("ignorekit"), bold("Commands:"), bold("Global flags:"), bold("Examples:"))
}
