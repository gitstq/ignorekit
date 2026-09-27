package cmd

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/gitstq/ignorekit/internal/engine"
	"github.com/gitstq/ignorekit/internal/gitignore"
	"github.com/gitstq/ignorekit/internal/gitx"
	"github.com/gitstq/ignorekit/internal/term"
)

func init() {
	register(command{name: "apply", short: "merge rules and optionally untrack", run: runApply})
}

func runApply(args []string) int {
	fs := flag.NewFlagSet("apply", flag.ContinueOnError)
	o := &commonOpts{}
	bindCommon(fs, o)
	dryRun := fs.Bool("dry-run", false, "show the plan without writing")
	if err := parse(fs, args); err != nil {
		return exitError
	}
	applyColor(o)

	res, err := engine.Analyze(o.dir)
	if err != nil {
		errf("error: %v\n", err)
		return exitError
	}

	merged := gitignore.Merge(res.GitignoreContent, res.Desired)
	gi := filepath.Join(res.Root, ".gitignore")

	if *dryRun {
		outf(term.Bold("📋 apply plan (dry-run)")+"\n")
		outf("  %s %d new rule(s) to add\n", term.Cyan("•"), len(merged.Added))
		for _, r := range merged.Added {
			outf("      + %s\n", term.Green(r))
		}
		if o.untrack {
			outf("  %s %d tracked file(s) to untrack (git rm --cached)\n", term.Cyan("•"), len(res.TrackedIssues))
			for i, iss := range res.TrackedIssues {
				if i >= maxList {
					outf("      … and %d more\n", len(res.TrackedIssues)-maxList)
					break
				}
				outf("      - %s\n", term.Yellow(iss.Path))
			}
		}
		return exitOK
	}

	if err := os.WriteFile(gi, []byte(merged.Content), 0o644); err != nil {
		errf("error writing .gitignore: %v\n", err)
		return exitError
	}
	outf("%s %s (%d rules in managed block, %d newly added)\n",
		term.Green("✔"), gi, len(merged.InManaged), len(merged.Added))

	if o.untrack && len(res.TrackedIssues) > 0 {
		if !res.IsGit {
			errf("%s not a git repository; skipped untrack\n", term.Yellow("!"))
		} else {
			paths := make([]string, 0, len(res.TrackedIssues))
			for _, iss := range res.TrackedIssues {
				paths = append(paths, iss.Path)
			}
			if err := gitx.Untrack(res.Root, paths); err != nil {
				errf("%s untrack failed: %v\n", term.Red("✗"), err)
				return exitError
			}
			outf("%s untracked %d file(s) (working-tree files kept)\n", term.Green("✔"), len(paths))
		}
	}

	if len(res.SecretIssues) > 0 {
		outf("%s %d secret red flag(s) in tracked files — rotate credentials!\n", term.Red("!"), len(res.SecretIssues))
	}
	return exitOK
}
