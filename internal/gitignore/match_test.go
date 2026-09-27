package gitignore

import "testing"

func dirs(set ...string) IsDir {
	m := map[string]bool{}
	for _, s := range set {
		m[s] = true
	}
	return func(rel string) bool { return m[rel] }
}

func TestCompileMatch(t *testing.T) {
	cases := []struct {
		pattern string
		path    string
		isDir   IsDir
		want    bool
	}{
		// Basic glob, unanchored.
		{"*.log", "error.log", nil, true},
		{"*.log", "sub/error.log", nil, true},
		{"*.log", "error.txt", nil, false},
		// Bare name matches at any depth and as a directory prefix.
		{"foo", "foo", nil, true},
		{"foo", "foo/bar", nil, true},
		{"foo", "a/foo", nil, true},
		{"foo", "a/foo/b", nil, true},
		{"foo", "foobar", nil, false},
		// Anchored by leading slash.
		{"/foo", "foo", nil, true},
		{"/foo", "foo/bar", nil, true},
		{"/foo", "a/foo", nil, false},
		// Directory-only pattern.
		{"foo/", "foo", dirs("foo"), true},
		{"foo/", "foo", nil, false},
		{"foo/", "foo/bar", nil, true},
		{"foo/", "a/foo/bar", nil, true},
		{"foo/", "a/foo", dirs("a/foo"), true},
		// Double star.
		{"**/foo", "foo", nil, true},
		{"**/foo", "a/foo", nil, true},
		{"**/foo", "a/b/foo", nil, true},
		{"abc/**", "abc/x", nil, true},
		{"abc/**", "abc/x/y", nil, true},
		{"abc/**", "abc", nil, false},
		{"a/**/b", "a/b", nil, true},
		{"a/x/b", "a/x/b", nil, true},
		{"a/**/b", "a/x/y/b", nil, true},
		{"a/**/b", "a/x", nil, false},
		// Single-char wildcard.
		{"foo?", "foox", nil, true},
		{"foo?", "foo", nil, false},
		{"foo?", "fooxy", nil, false},
		// Character classes.
		{"[abc]*", "apple", nil, true},
		{"[abc]*", "xyz", nil, false},
		{"[!abc]*", "xyz", nil, true},
		{"[!abc]*", "apple", nil, false},
		// Anchored nested path.
		{"bar/foo", "bar/foo", nil, true},
		{"bar/foo", "bar/foo/baz", nil, true},
		{"bar/foo", "x/bar/foo", nil, false},
	}
	for _, c := range cases {
		p, err := Compile(c.pattern, "")
		if err != nil {
			t.Fatalf("compile %q: %v", c.pattern, err)
		}
		if p == nil {
			t.Fatalf("pattern %q produced nil", c.pattern)
		}
		got := p.Match(c.path, c.isDir)
		if got != c.want {
			t.Errorf("pattern %q vs path %q: got %v want %v", c.pattern, c.path, got, c.want)
		}
	}
}

func TestCompileBlankAndComment(t *testing.T) {
	for _, line := range []string{"", "   ", "# just a comment", "\t\n"} {
		p, err := Compile(line, "")
		if err != nil {
			t.Fatalf("compile %q: %v", line, err)
		}
		if p != nil {
			t.Errorf("line %q should produce nil", line)
		}
	}
}

func TestFileSetNegation(t *testing.T) {
	fs := NewFileSet()
	content := "*.log\n!important.log\n"
	if err := fs.ParseReader(content, ""); err != nil {
		t.Fatal(err)
	}
	if !fs.Ignored("debug.log", nil) {
		t.Error("debug.log should be ignored")
	}
	if fs.Ignored("important.log", nil) {
		t.Error("important.log should be re-included by negation")
	}
}

func TestFileSetLastMatchWins(t *testing.T) {
	fs := NewFileSet()
	// foo ignored, then the whole thing re-included by a later broad negation.
	if err := fs.ParseReader("foo\n!foo\n", ""); err != nil {
		t.Fatal(err)
	}
	if fs.Ignored("foo", nil) {
		t.Error("later negation should win")
	}
}
