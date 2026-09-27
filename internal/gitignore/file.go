package gitignore

import (
	"bufio"
	"strings"
)

// FileSet is an ordered collection of patterns parsed from one or more
// gitignore sources. Git applies rules with "last matching pattern wins"
// semantics, including negation.
type FileSet struct {
	Patterns []*Pattern
}

// NewFileSet returns an empty FileSet.
func NewFileSet() *FileSet { return &FileSet{} }

// ParseReader appends rules from line-oriented text content.
func (f *FileSet) ParseReader(content, section string) error {
	sc := bufio.NewScanner(strings.NewReader(content))
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		p, err := Compile(line, section)
		if err != nil {
			return err
		}
		if p != nil {
			f.Patterns = append(f.Patterns, p)
		}
	}
	return sc.Err()
}

// Ignored reports whether the relative path is ignored by this rule set.
// It walks every pattern in order; the last matching pattern decides.
// A path with no matching pattern is not ignored.
func (f *FileSet) Ignored(rel string, isDir IsDir) bool {
	ignored := false
	for _, p := range f.Patterns {
		if p.Match(rel, isDir) {
			ignored = !p.Negate
		}
	}
	return ignored
}

// MatchPattern returns the last pattern that matches rel, or nil.
func (f *FileSet) MatchPattern(rel string, isDir IsDir) *Pattern {
	var hit *Pattern
	for _, p := range f.Patterns {
		if p.Match(rel, isDir) {
			hit = p
		}
	}
	return hit
}
