package gitignore

import (
	"strings"
	"testing"
)

func TestSplitManaged(t *testing.T) {
	existing := "user-rule\n\n" + BeginMarker + "\n# --- X ---\nold\n" + EndMarker + "\n# tail"
	user, managed, found := SplitManaged(existing)
	if !found {
		t.Fatal("managed block should be detected")
	}
	if !strings.Contains(user, "user-rule") {
		t.Errorf("user content missing: %q", user)
	}
	if strings.Contains(user, BeginMarker) {
		t.Error("markers should be removed from user content")
	}
	if !strings.Contains(managed, "old") {
		t.Errorf("managed body missing: %q", managed)
	}
	if !strings.Contains(user, "# tail") {
		t.Errorf("trailing user content should be preserved: %q", user)
	}
}

func TestMergeAddsAndDedupes(t *testing.T) {
	desired := []Rule{
		{Line: "node_modules/", Section: "Node.js"},
		{Line: "*.log", Section: "Logs"},
	}
	res := Merge("", desired)
	if len(res.Added) != 2 {
		t.Fatalf("expected 2 added, got %d (%v)", len(res.Added), res.Added)
	}
	if !strings.Contains(res.Content, BeginMarker) || !strings.Contains(res.Content, EndMarker) {
		t.Error("managed markers missing")
	}
	if !strings.Contains(res.Content, "node_modules/") || !strings.Contains(res.Content, "*.log") {
		t.Error("rules missing from merged content")
	}

	// Re-running with the same rules must not report new additions.
	again := Merge(res.Content, desired)
	if len(again.Added) != 0 {
		t.Errorf("second run should add nothing, got %v", again.Added)
	}
}

func TestMergePreservesHandRules(t *testing.T) {
	existing := "# my hand-written rule\nmy-stuff/\n"
	desired := []Rule{{Line: "node_modules/", Section: "Node.js"}}
	res := Merge(existing, desired)
	if !strings.Contains(res.Content, "my-stuff/") {
		t.Error("hand-written rule must be preserved")
	}
	if !strings.Contains(res.Content, "node_modules/") {
		t.Error("generated rule must be added")
	}

	// A rule already declared by hand should not be duplicated in the block.
	desired2 := []Rule{{Line: "my-stuff/", Section: "Mine"}}
	res2 := Merge(existing, desired2)
	if strings.Count(res2.Content, "my-stuff/") != 1 {
		t.Errorf("hand rule should appear exactly once:\n%s", res2.Content)
	}
}

func TestMergeNewRulesDetected(t *testing.T) {
	existing := BeginMarker + "\nnode_modules/\n" + EndMarker + "\n"
	desired := []Rule{
		{Line: "node_modules/", Section: "Node.js"},
		{Line: "*.log", Section: "Logs"},
	}
	res := Merge(existing, desired)
	if len(res.Added) != 1 || res.Added[0] != "*.log" {
		t.Errorf("only *.log should be new, got %v", res.Added)
	}
}
