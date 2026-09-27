package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gitstq/ignorekit/internal/gitx"
)

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func write(t *testing.T, dir, rel, content string) {
	t.Helper()
	full := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestEndToEndTrackingHygiene(t *testing.T) {
	if !gitx.Available() {
		t.Skip("git not available")
	}
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	git(t, dir, "config", "commit.gpgsign", "false")

	write(t, dir, "package.json", "{}")
	write(t, dir, "node_modules/lib/index.js", "lib")
	write(t, dir, ".env", "TOKEN=ghp_0123456789abcdef0123456789abcdef0123")
	write(t, dir, ".claude/settings.local.json", "{}")
	write(t, dir, "src/main.js", "console.log('ok')")

	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", "init")

	res, err := Analyze(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsGit {
		t.Fatal("should be detected as a git repo")
	}

	flagged := map[string]bool{}
	for _, iss := range res.TrackedIssues {
		flagged[iss.Path] = true
	}
	for _, want := range []string{"node_modules/lib/index.js", ".env", ".claude/settings.local.json"} {
		if !flagged[want] {
			t.Errorf("expected %q to be flagged as tracked-but-ignored", want)
		}
	}
	if flagged["src/main.js"] {
		t.Error("src/main.js should not be flagged")
	}

	if len(res.SecretIssues) == 0 {
		t.Error("expected a secret red flag in .env")
	} else {
		found := false
		for _, s := range res.SecretIssues {
			if s.File == ".env" {
				found = true
			}
		}
		if !found {
			t.Errorf("secret should be located in .env, got %+v", res.SecretIssues)
		}
	}

	// Now untrack the flagged files and verify they leave the index.
	paths := make([]string, 0, len(res.TrackedIssues))
	for _, iss := range res.TrackedIssues {
		paths = append(paths, iss.Path)
	}
	if err := gitx.Untrack(dir, paths); err != nil {
		t.Fatal(err)
	}
	tracked, err := gitx.TrackedFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	still := map[string]bool{}
	for _, f := range tracked {
		still[f] = true
	}
	for _, p := range paths {
		if still[p] {
			t.Errorf("%q should no longer be tracked", p)
		}
	}
	if !still["src/main.js"] {
		t.Error("src/main.js should remain tracked")
	}
}
