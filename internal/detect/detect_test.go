package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func mkdirp(t *testing.T, root string, paths ...string) {
	t.Helper()
	for _, p := range paths {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(full, 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func touch(t *testing.T, root string, paths ...string) {
	t.Helper()
	for _, p := range paths {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func hasKey(ds []Detection, key string) bool {
	for _, d := range ds {
		if d.Key == key {
			return true
		}
	}
	return false
}

func TestScanNodeProject(t *testing.T) {
	root := t.TempDir()
	touch(t, root, "package.json", "package-lock.json")
	mkdirp(t, root, "node_modules")
	ds, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if !hasKey(ds, "node") {
		t.Errorf("node should be detected: %+v", ds)
	}
}

func TestScanPythonAndGo(t *testing.T) {
	root := t.TempDir()
	touch(t, root, "pyproject.toml", "go.mod")
	ds, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if !hasKey(ds, "python") || !hasKey(ds, "go") {
		t.Errorf("python and go should be detected: %+v", ds)
	}
}

func TestScanAIAgents(t *testing.T) {
	root := t.TempDir()
	touch(t, root, ".claude/settings.local.json", ".cursor/mcp.json", ".aider.input.history")
	mkdirp(t, root, ".claude", ".cursor")
	ds, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"claude-code", "cursor", "aider"} {
		if !hasKey(ds, key) {
			t.Errorf("%s should be detected: %+v", key, ds)
		}
	}
}

func TestScanEmpty(t *testing.T) {
	root := t.TempDir()
	ds, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) != 0 {
		t.Errorf("empty dir should detect nothing, got %+v", ds)
	}
}

func TestGlobMarker(t *testing.T) {
	root := t.TempDir()
	touch(t, root, "MyApp.sln", "MyApp.csproj")
	ds, err := Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	if !hasKey(ds, "dotnet") {
		t.Errorf("dotnet should be detected via glob: %+v", ds)
	}
}
