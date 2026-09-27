// Package gitx wraps the small set of git operations ignorekit needs.
// It shells out to the git binary; matching logic itself is pure Go.
package gitx

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
)

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String(), &GitError{Args: args, Stderr: strings.TrimSpace(errb.String()), Err: err}
	}
	return out.String(), nil
}

// GitError carries git's stderr for clearer messages.
type GitError struct {
	Args   []string
	Stderr string
	Err    error
}

func (e *GitError) Error() string {
	if e.Stderr != "" {
		return "git " + strings.Join(e.Args, " ") + ": " + e.Stderr
	}
	return e.Err.Error()
}

// Available reports whether the git binary is installed.
func Available() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// IsRepo reports whether dir is inside a git working tree.
func IsRepo(dir string) bool {
	out, err := git(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// Root returns the repository top-level directory.
func Root(dir string) (string, error) {
	out, err := git(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// TrackedFiles lists every path currently tracked in the index/HEAD,
// returned as slash-separated paths relative to the repository root.
func TrackedFiles(dir string) ([]string, error) {
	out, err := git(dir, "ls-files", "--cached", "--exclude-standard")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, l := range strings.Split(out, "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			files = append(files, filepath.ToSlash(l))
		}
	}
	return files, nil
}

// Untrack removes the given paths from the index (keeping working-tree files),
// batching to avoid overly long command lines.
func Untrack(dir string, files []string) error {
	const batch = 200
	for i := 0; i < len(files); i += batch {
		j := i + batch
		if j > len(files) {
			j = len(files)
		}
		args := []string{"rm", "--cached", "--quiet", "--"}
		args = append(args, files[i:j]...)
		if _, err := git(dir, args...); err != nil {
			return err
		}
	}
	return nil
}
