// Package secret provides a small, high-confidence set of "red flag" patterns
// for obvious credentials. It is intentionally not a full secret scanner;
// ignorekit points users to dedicated tools for comprehensive coverage.
package secret

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// maxScanBytes limits how much of a file is inspected.
const maxScanBytes = 1 << 20 // 1 MiB

// Rule is one high-confidence credential pattern.
type Rule struct {
	Type string
	RE   *regexp.Regexp
}

// redact masks all but a short prefix of a matched secret.
func redact(match string) string {
	if len(match) <= 12 {
		return match[:min(4, len(match))] + "***"
	}
	return match[:8] + "***" + match[len(match)-2:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Rules are deliberately high-signal to limit false positives.
var Rules = []Rule{
	{"GitHub personal access token", regexp.MustCompile(`\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{36,}\b`)},
	{"GitHub fine-grained token", regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{60,}\b`)},
	{"GitLab personal access token", regexp.MustCompile(`\bglpat-[A-Za-z0-9_\-]{20,}\b`)},
	{"AWS access key ID", regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`)},
	{"Google API key", regexp.MustCompile(`\bAIza[0-9A-Za-z_\-]{35}\b`)},
	{"Slack token", regexp.MustCompile(`\bxox[baprs]-[0-9A-Za-z-]{10,}\b`)},
	{"Stripe live secret key", regexp.MustCompile(`\b(?:sk|rk)_live_[0-9A-Za-z]{20,}\b`)},
	{"OpenAI API key", regexp.MustCompile(`\bsk-(?:proj-)?[A-Za-z0-9]{20,}\b`)},
	{"JWT", regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{5,}\.eyJ[A-Za-z0-9_-]{5,}\.[A-Za-z0-9_-]{5,}\b`)},
	{"PEM private key", regexp.MustCompile(`-----BEGIN (?:[A-Z ]+ )?PRIVATE KEY-----`)},
}

// Finding is one red flag located in a file.
type Finding struct {
	File     string
	Type     string
	Line     int
	Redacted string
}

// ScanFile inspects a single file.
func ScanFile(path string) ([]Finding, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var findings []Finding
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), maxScanBytes)
	line := 0
	for sc.Scan() {
		line++
		text := sc.Text()
		for _, r := range Rules {
			if r.RE.MatchString(text) {
				m := r.RE.FindString(text)
				findings = append(findings, Finding{
					File:     path,
					Type:     r.Type,
					Line:     line,
					Redacted: redact(strings.TrimSpace(m)),
				})
			}
		}
	}
	return findings, sc.Err()
}

// ScanPaths inspects the given root-relative paths under root.
func ScanPaths(root string, rels []string) ([]Finding, error) {
	var all []Finding
	for _, rel := range rels {
		full := filepath.Join(root, filepath.FromSlash(rel))
		fi, err := os.Stat(full)
		if err != nil || fi.IsDir() || fi.Size() > maxScanBytes {
			continue
		}
		fs, err := ScanFile(full)
		if err != nil {
			continue
		}
		for i := range fs {
			fs[i].File = rel
		}
		all = append(all, fs...)
	}
	return all, nil
}
