package cmd

import (
	"flag"

	"github.com/gitstq/ignorekit/internal/engine"
	"github.com/gitstq/ignorekit/internal/report"
)

func init() {
	register(command{name: "audit", short: "read-only hygiene report", run: runAudit})
}

func runAudit(args []string) int {
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	o := &commonOpts{}
	bindCommon(fs, o)
	strict := fs.Bool("strict", false, "treat suggested-rule additions as findings (CI gate)")
	if err := parse(fs, args); err != nil {
		return exitError
	}
	applyColor(o)

	res, err := engine.Analyze(o.dir)
	if err != nil {
		errf("error: %v\n", err)
		return exitError
	}

	switch o.format {
	case "json":
		b, err := report.JSON(res)
		if err != nil {
			errf("error: %v\n", err)
			return exitError
		}
		outf("%s\n", b)
	case "sarif":
		b, err := report.SARIF(res)
		if err != nil {
			errf("error: %v\n", err)
			return exitError
		}
		outf("%s\n", b)
	case "text", "":
		outf("%s", report.Text(res))
	default:
		errf("unknown format %q (use text|json|sarif)\n", o.format)
		return exitError
	}

	if len(res.SecretIssues) > 0 || len(res.TrackedIssues) > 0 {
		return exitFindings
	}
	if *strict && len(res.Desired) > 0 && !res.GitignorePresent {
		return exitFindings
	}
	return exitOK
}
