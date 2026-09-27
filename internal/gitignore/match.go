// Package gitignore parses .gitignore files, compiles their patterns,
// matches relative paths, and safely merges generated rules with existing ones.
// It is a self-contained implementation using only the Go standard library.
package gitignore

import (
	"regexp"
	"strings"
)

// Pattern is a single compiled gitignore rule.
type Pattern struct {
	Raw      string // original line (without trailing newline)
	Negate   bool   // line started with '!'
	DirOnly  bool   // line ended with '/'
	Anchored bool   // pattern is anchored to the ignore-file root
	Section  string // human-readable origin of the rule (for merge grouping)

	// sig is a normalized identity used for de-duplication.
	sig string
	// matchRE matches a file path (and, for non-dir-only rules, the bare name).
	matchRE *regexp.Regexp
	// bareRE is set only for DirOnly rules; it matches the bare directory path.
	bareRE *regexp.Regexp
}

// IsDir reports whether a slash-separated relative path is a directory.
type IsDir func(rel string) bool

// Match reports whether the pattern matches the given relative path.
// rel must use forward slashes and be relative to the ignore-file root.
// isDir is consulted only for directory-only patterns matching a bare path;
// it may be nil (bare directory paths then never match).
func (p *Pattern) Match(rel string, isDir IsDir) bool {
	if p.DirOnly {
		if p.matchRE != nil && p.matchRE.MatchString(rel) {
			return true
		}
		if p.bareRE != nil && p.bareRE.MatchString(rel) {
			if isDir != nil && isDir(rel) {
				return true
			}
		}
		return false
	}
	return p.matchRE != nil && p.matchRE.MatchString(rel)
}

var (
	extraStar = regexp.MustCompile(`\*+`)
)

// translateSegment converts one path segment (no slashes) to a regexp fragment.
func translateSegment(seg string) string {
	var b strings.Builder
	i := 0
	for i < len(seg) {
		c := seg[i]
		switch {
		case c == '\\' && i+1 < len(seg):
			// Escaped literal character.
			b.WriteString(regexp.QuoteMeta(string(seg[i+1])))
			i += 2
		case c == '*':
			if i+1 < len(seg) && seg[i+1] == '*' {
				b.WriteString(".*")
				i += 2
				for i < len(seg) && seg[i] == '*' {
					i++
				}
			} else {
				b.WriteString("[^/]*")
				i++
			}
		case c == '?':
			b.WriteString("[^/]")
			i++
		case c == '[':
			j := i + 1
			if j < len(seg) && (seg[j] == '!' || seg[j] == '^') {
				j++
			}
			if j < len(seg) && seg[j] == ']' {
				j++
			}
			for j < len(seg) && seg[j] != ']' {
				j++
			}
			if j >= len(seg) {
				b.WriteString(`\[`)
				i++
			} else {
				cls := seg[i : j+1]
				if len(cls) > 1 && cls[1] == '!' {
					cls = "[^" + cls[2:]
				}
				b.WriteString(cls)
				i = j + 1
			}
		default:
			b.WriteString(regexp.QuoteMeta(string(c)))
			i++
		}
	}
	return b.String()
}

// translateBody converts a pattern (anchoring/trailing slash already stripped)
// into a regexp fragment, handling the cross-directory '**' segments.
func translateBody(pat string) string {
	segs := strings.Split(pat, "/")
	var b strings.Builder
	for idx, seg := range segs {
		if seg == "**" {
			switch {
			case idx == 0 && len(segs) == 1:
				b.WriteString(".*")
			case idx == 0:
				b.WriteString("(?:.*/)?") // leading **/
			case idx == len(segs)-1:
				b.WriteString("/.*") // trailing /**
			default:
				// middle **: one connecting slash, then zero or more dirs
				// each ending in a slash; the following segment adds no slash.
				b.WriteString("/(?:(?:[^/]+/)*)?")
			}
			continue
		}
		if idx > 0 && segs[idx-1] != "**" {
			b.WriteString("/")
		}
		b.WriteString(translateSegment(seg))
	}
	return b.String()
}

// Compile parses a single gitignore rule line. Blank lines and comment lines
// return nil with no error.
func Compile(line, section string) (*Pattern, error) {
	raw := line
	trimmed := strings.TrimRight(line, "\r\n")
	if strings.TrimSpace(trimmed) == "" {
		return nil, nil
	}
	if strings.HasPrefix(trimmed, "#") {
		return nil, nil
	}

	p := &Pattern{Raw: raw, Section: section}
	body := trimmed
	if strings.HasPrefix(body, "!") {
		p.Negate = true
		body = body[1:]
	}
	if strings.HasSuffix(body, "/") {
		p.DirOnly = true
		body = body[:len(body)-1]
	}
	if strings.HasPrefix(body, "/") {
		body = body[1:]
		p.Anchored = true
	} else if strings.Contains(body, "/") {
		p.Anchored = true
	}

	frag := translateBody(body)
	prefix := "(?:.*/)?"
	if p.Anchored {
		prefix = ""
	}

	{
		var s strings.Builder
		if p.Negate {
			s.WriteByte('!')
		}
		if p.Anchored {
			s.WriteByte('/')
		}
		s.WriteString(body)
		if p.DirOnly {
			s.WriteByte('/')
		}
		p.sig = s.String()
	}

	if p.DirOnly {
		re, err := regexp.Compile("^" + prefix + frag + "/.*$")
		if err != nil {
			return nil, err
		}
		bare, err := regexp.Compile("^" + prefix + frag + "$")
		if err != nil {
			return nil, err
		}
		p.matchRE = re
		p.bareRE = bare
		return p, nil
	}

	var expr string
	if strings.HasSuffix(frag, "/.*") {
		expr = "^" + prefix + frag + "$"
	} else {
		expr = "^" + prefix + frag + "(?:/.*)?$"
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	p.matchRE = re
	return p, nil
}
