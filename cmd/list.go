package cmd

import (
	"flag"
	"sort"

	"github.com/gitstq/ignorekit/internal/templates"
	"github.com/gitstq/ignorekit/internal/term"
)

func init() {
	register(command{name: "list-templates", short: "list built-in templates", run: runListTemplates})
}

func runListTemplates(args []string) int {
	fs := flag.NewFlagSet("list-templates", flag.ContinueOnError)
	o := &commonOpts{}
	bindCommon(fs, o)
	if err := parse(fs, args); err != nil {
		return exitError
	}
	applyColor(o)

	all := templates.All()
	sort.Slice(all, func(i, j int) bool {
		if all[i].Order != all[j].Order {
			return all[i].Order < all[j].Order
		}
		return all[i].Key < all[j].Key
	})

	lastCat := ""
	for _, t := range all {
		if t.Category != lastCat {
			outf("\n%s\n", term.Purple("["+t.Category+"]"))
			lastCat = t.Category
		}
		outf("  %-14s %s  %s\n", term.Bold(t.Key), term.Cyan(t.Title), term.Gray(plural(len(t.Rules), "rule")))
	}
	return exitOK
}

func plural(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return itoa(n) + " " + word + "s"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
