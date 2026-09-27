// Package engine orchestrates detection, template assembly, gitignore
// matching against tracked files, and secret red flags into a single Result.
package engine

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gitstq/ignorekit/internal/detect"
	"github.com/gitstq/ignorekit/internal/gitignore"
	"github.com/gitstq/ignorekit/internal/gitx"
	"github.com/gitstq/ignorekit/internal/secret"
	"github.com/gitstq/ignorekit/internal/templates"
)

// ToolInfo summarizes a detected tool.
type ToolInfo struct {
	Key      string
	Name     string
	Category string
	Matched  []string
	Shared   []string
	Review   []string
	Notes    string
}

// TrackedIssue is a tracked file that should be ignored.
type TrackedIssue struct {
	Path        string
	MatchedRule string
	Section     string
	FromNewRule bool
}

// Result is the full analysis outcome.
type Result struct {
	Root             string
	IsGit            bool
	Tools            []ToolInfo
	Desired          []gitignore.Rule
	GitignorePresent bool
	GitignoreContent string
	TrackedIssues    []TrackedIssue
	SecretIssues     []secret.Finding
	TrackedTotal     int
}

// Analyze scans root and builds a Result.
func Analyze(root string) (*Result, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	res := &Result{Root: abs, IsGit: gitx.Available() && gitx.IsRepo(abs)}

	detections, err := detect.Scan(abs)
	if err != nil {
		return nil, err
	}

	// Map detected keys and preserve template ordering.
	keySet := map[string]detect.Detection{}
	for _, d := range detections {
		res.Tools = append(res.Tools, ToolInfo{
			Key: d.Key, Name: d.Name, Category: d.Category,
			Matched: d.Matched, Shared: d.Shared, Review: d.Review, Notes: d.Notes,
		})
		keySet[d.Key] = d
	}

	var keys []string
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ti, _ := templates.Get(keys[i])
		tj, _ := templates.Get(keys[j])
		if ti.Order != tj.Order {
			return ti.Order < tj.Order
		}
		return ti.Key < tj.Key
	})

	for _, k := range keys {
		t, ok := templates.Get(k)
		if !ok {
			continue
		}
		for _, r := range t.Rules {
			res.Desired = append(res.Desired, gitignore.Rule{Line: r, Section: t.Title})
		}
	}

	// Read existing .gitignore.
	giPath := filepath.Join(abs, ".gitignore")
	if b, rerr := os.ReadFile(giPath); rerr == nil {
		res.GitignorePresent = true
		res.GitignoreContent = string(b)
	}

	if !res.IsGit {
		return res, nil
	}

	tracked, err := gitx.TrackedFiles(abs)
	if err != nil {
		return res, nil
	}
	res.TrackedTotal = len(tracked)

	existingSet := gitignore.NewFileSet()
	_ = existingSet.ParseReader(res.GitignoreContent, "existing")
	desiredSet := gitignore.NewFileSet()
	for _, r := range res.Desired {
		_ = desiredSet.ParseReader(r.Line, r.Section)
	}
	combinedSet := gitignore.NewFileSet()
	_ = combinedSet.ParseReader(res.GitignoreContent, "existing")
	for _, r := range res.Desired {
		_ = combinedSet.ParseReader(r.Line, r.Section)
	}

	isDir := func(rel string) bool {
		fi, statErr := os.Stat(filepath.Join(abs, filepath.FromSlash(rel)))
		return statErr == nil && fi.IsDir()
	}

	for _, f := range tracked {
		// A tracked file that the combined rules consider ignored is a hygiene issue.
		if p := combinedSet.MatchPattern(f, isDir); p != nil {
			issue := TrackedIssue{
				Path:        f,
				MatchedRule: p.Raw,
				Section:     p.Section,
				FromNewRule: desiredSet.MatchPattern(f, isDir) != nil && existingSet.MatchPattern(f, isDir) == nil,
			}
			if issue.Section == "" || issue.Section == "existing" {
				issue.Section = "existing .gitignore"
			}
			res.TrackedIssues = append(res.TrackedIssues, issue)
		}
	}

	// Secret red flags: scan tracked files plus any MCP/env config present.
	scanRels := append([]string{}, tracked...)
	for _, extra := range []string{
		".mcp.json", "mcp.json", ".cursor/mcp.json", ".vscode/mcp.json",
		".claude/settings.json", ".claude/settings.local.json",
	} {
		if _, statErr := os.Stat(filepath.Join(abs, filepath.FromSlash(extra))); statErr == nil {
			scanRels = append(scanRels, extra)
		}
	}
	seen := map[string]bool{}
	uniq := scanRels[:0]
	for _, r := range scanRels {
		if !seen[r] {
			seen[r] = true
			uniq = append(uniq, r)
		}
	}
	if sf, err := secret.ScanPaths(abs, uniq); err == nil {
		res.SecretIssues = sf
	}

	return res, nil
}
