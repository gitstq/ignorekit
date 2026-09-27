// Package report renders an engine.Result as human-readable text, JSON, or
// SARIF 2.1.0 for CI integration.
package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gitstq/ignorekit/internal/engine"
	"github.com/gitstq/ignorekit/internal/secret"
	"github.com/gitstq/ignorekit/internal/term"
)

const maxList = 15

// Text renders a colored human-readable report.
func Text(res *engine.Result) string {
	var b strings.Builder

	b.WriteString(term.Bold(term.Cyan("🦞 ignorekit · tracking-hygiene audit")) + "\n")
	b.WriteString(term.Gray("root: "+res.Root) + "\n\n")

	// Detected tools grouped by category.
	b.WriteString(term.Bold("🔎 Detected tools") + "\n")
	if len(res.Tools) == 0 {
		b.WriteString(term.Yellow("   (no recognizable project markers)") + "\n")
	} else {
		catOrder := []string{"Language", "Framework", "Editor", "AI", "Secrets", "Deploy", "Logs"}
		catSeen := map[string]bool{}
		for _, cat := range catOrder {
			var group []engine.ToolInfo
			for _, t := range res.Tools {
				if t.Category == cat {
					group = append(group, t)
				}
			}
			if len(group) == 0 {
				continue
			}
			catSeen[cat] = true
			b.WriteString("   " + term.Purple("["+cat+"]") + "\n")
			for _, t := range group {
				b.WriteString("    • " + term.Bold(t.Name) + "\n")
				for _, s := range t.Shared {
					b.WriteString("        " + term.Green("✔ keep/commit: ") + term.Green(s) + "\n")
				}
				for _, rv := range t.Review {
					b.WriteString("        " + term.Yellow("⚠ review for secrets: ") + term.Yellow(rv) + "\n")
				}
				if t.Notes != "" {
					b.WriteString("        " + term.Gray(t.Notes) + "\n")
				}
			}
		}
		for _, t := range res.Tools {
			if !catSeen[t.Category] {
				b.WriteString("    • " + term.Bold(t.Name) + term.Gray("  ["+t.Category+"]") + "\n")
			}
		}
	}
	b.WriteString("\n")

	// .gitignore status.
	b.WriteString(term.Bold("📄 .gitignore") + "\n")
	if res.GitignorePresent {
		b.WriteString("   " + term.Green("present") + "\n")
	} else {
		b.WriteString("   " + term.Yellow("missing") + "\n")
	}
	b.WriteString("   " + term.Gray(fmt.Sprintf("%d rules suggested", len(res.Desired))) + "\n\n")

	// Tracked issues.
	b.WriteString(term.Bold("🧹 Tracked files that should be ignored") + "\n")
	if !res.IsGit {
		b.WriteString("   " + term.Gray("(not a git repository — skipped)") + "\n\n")
	} else if len(res.TrackedIssues) == 0 {
		b.WriteString("   " + term.Green("none — clean ✨") + "\n\n")
	} else {
		b.WriteString("   " + term.Yellow(fmt.Sprintf("%d file(s)", len(res.TrackedIssues))) + "\n")
		list := res.TrackedIssues
		sort.SliceStable(list, func(i, j int) bool {
			if list[i].FromNewRule != list[j].FromNewRule {
				return list[i].FromNewRule
			}
			return list[i].Path < list[j].Path
		})
		for i, iss := range list {
			if i >= maxList {
				b.WriteString("   " + term.Gray(fmt.Sprintf("… and %d more", len(list)-maxList)) + "\n")
				break
			}
			tag := term.Gray("existing rule")
			if iss.FromNewRule {
				tag = term.Yellow("new rule")
			}
			b.WriteString(fmt.Sprintf("    %s  %s  %s\n", term.Red("✗"), iss.Path, term.Gray("("+iss.Section+" · "+tag+")")))
		}
		b.WriteString("\n")
	}

	// Secret issues.
	b.WriteString(term.Bold("🔐 Secret red flags") + "\n")
	if len(res.SecretIssues) == 0 {
		b.WriteString("   " + term.Green("no high-confidence secrets found") + "\n\n")
	} else {
		b.WriteString("   " + term.Red(fmt.Sprintf("%d finding(s) — rotate these credentials now!", len(res.SecretIssues))) + "\n")
		for i, s := range res.SecretIssues {
			if i >= maxList {
				b.WriteString("   " + term.Gray(fmt.Sprintf("… and %d more", len(res.SecretIssues)-maxList)) + "\n")
				break
			}
			b.WriteString(fmt.Sprintf("    %s %s:%d  %s  %s\n", term.Red("✗"), s.File, s.Line, term.Red(s.Type), term.Gray(s.Redacted)))
		}
		b.WriteString("\n")
	}

	// Next steps.
	b.WriteString(term.Bold("👉 Next steps") + "\n")
	if !res.GitignorePresent {
		b.WriteString("   • " + term.Cyan("ignorekit init --write") + term.Gray("   create a .gitignore\n"))
	} else {
		b.WriteString("   • " + term.Cyan("ignorekit apply") + term.Gray("        merge suggested rules safely\n"))
	}
	if len(res.TrackedIssues) > 0 {
		b.WriteString("   • " + term.Cyan("ignorekit apply --untrack") + term.Gray("  also `git rm --cached` these files\n"))
	}
	if len(res.SecretIssues) > 0 {
		b.WriteString("   • " + term.Red("rotate exposed credentials, then commit the cleanup\n"))
	}
	return b.String()
}

type jsonTool struct {
	Key      string   `json:"key"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Shared   []string `json:"shared,omitempty"`
	Review   []string `json:"review,omitempty"`
}

type jsonReport struct {
	Root             string                `json:"root"`
	IsGit             bool                  `json:"is_git"`
	GitignorePresent  bool                  `json:"gitignore_present"`
	Tools             []jsonTool            `json:"tools"`
	DesiredRules      []string              `json:"desired_rules"`
	TrackedTotal      int                   `json:"tracked_total"`
	TrackedIssues     []engine.TrackedIssue `json:"tracked_issues"`
	SecretIssues      []secret.Finding      `json:"secret_issues"`
	Summary           map[string]int        `json:"summary"`
}

// JSON renders a machine-readable report.
func JSON(res *engine.Result) ([]byte, error) {
	jr := jsonReport{
		Root: res.Root, IsGit: res.IsGit, GitignorePresent: res.GitignorePresent,
		TrackedTotal:  res.TrackedTotal,
		TrackedIssues: res.TrackedIssues,
		SecretIssues:  res.SecretIssues,
		Summary: map[string]int{
			"tools":           len(res.Tools),
			"rules_suggested": len(res.Desired),
			"tracked_issues":  len(res.TrackedIssues),
			"secret_issues":   len(res.SecretIssues),
		},
	}
	for _, t := range res.Tools {
		jr.Tools = append(jr.Tools, jsonTool{Key: t.Key, Name: t.Name, Category: t.Category, Shared: t.Shared, Review: t.Review})
	}
	for _, r := range res.Desired {
		jr.DesiredRules = append(jr.DesiredRules, r.Line)
	}
	// Guarantee arrays (never null) for stable CI consumption.
	if jr.Tools == nil {
		jr.Tools = []jsonTool{}
	}
	if jr.DesiredRules == nil {
		jr.DesiredRules = []string{}
	}
	if jr.TrackedIssues == nil {
		jr.TrackedIssues = []engine.TrackedIssue{}
	}
	if jr.SecretIssues == nil {
		jr.SecretIssues = []secret.Finding{}
	}
	return json.MarshalIndent(jr, "", "  ")
}

// SARIF renders SARIF 2.1.0 suitable for code-scanning CI.
func SARIF(res *engine.Result) ([]byte, error) {
	type region struct {
		StartLine int `json:"startLine"`
	}
	type location struct {
		PhysicalLocation struct {
			ArtifactLocation struct {
				URI string `json:"uri"`
			} `json:"artifactLocation"`
			Region *region `json:"region,omitempty"`
		} `json:"physicalLocation"`
	}
	type result struct {
		RuleID    string     `json:"ruleId"`
		Level     string     `json:"level"`
		Message   struct{ Text string } `json:"message"`
		Locations []location `json:"locations,omitempty"`
	}
	type rule struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		ShortDescription struct{ Text string } `json:"shortDescription"`
		DefaultConfiguration struct {
			Level string `json:"level"`
		} `json:"defaultConfiguration"`
	}

	mkRule := func(id, name, desc, level string) rule {
		r := rule{ID: id, Name: name}
		r.ShortDescription.Text = desc
		r.DefaultConfiguration.Level = level
		return r
	}
	rules := []rule{
		mkRule("IGK001", "tracked-ignored-file", "File is tracked by git but matches an ignore rule", "warning"),
		mkRule("IGK002", "possible-secret", "High-confidence credential or private key found", "error"),
	}

	var results = []result{}
	for _, iss := range res.TrackedIssues {
		r := result{RuleID: "IGK001", Level: "warning"}
		r.Message.Text = fmt.Sprintf("Tracked file %q matches ignore rule %q (%s)", iss.Path, iss.MatchedRule, iss.Section)
		loc := location{}
		loc.PhysicalLocation.ArtifactLocation.URI = iss.Path
		r.Locations = []location{loc}
		results = append(results, r)
	}
	for _, s := range res.SecretIssues {
		r := result{RuleID: "IGK002", Level: "error"}
		r.Message.Text = fmt.Sprintf("Possible %s in %s", s.Type, s.File)
		loc := location{}
		loc.PhysicalLocation.ArtifactLocation.URI = s.File
		loc.PhysicalLocation.Region = &region{StartLine: s.Line}
		r.Locations = []location{loc}
		results = append(results, r)
	}

	doc := map[string]interface{}{
		"$schema": "https://json.schemastore.org/sarif-2.1.0.json",
		"version": "2.1.0",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":            "ignorekit",
						"informationUri":  "https://github.com/gitstq/ignorekit",
						"version":         "1.0.0",
						"rules":           rules,
					},
				},
				"results": results,
			},
		},
	}
	return json.MarshalIndent(doc, "", "  ")
}
