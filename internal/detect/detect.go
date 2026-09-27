package detect

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Detection is a signature that matched the scanned project.
type Detection struct {
	Signature
	Matched []string // marker paths that were present
}

// heavyDirs are dependency/build directories we record but never descend into.
var heavyDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true, "target": true,
	"dist": true, "build": true, "out": true, ".next": true, ".nuxt": true,
	".svelte-kit": true, ".venv": true, "venv": true, "env": true,
	"__pycache__": true, ".gradle": true, ".dart_tool": true, ".expo": true,
	".cache": true, ".parcel-cache": true, ".turbo": true, "storybook-static": true,
}

func depth(rel string) int { return strings.Count(rel, "/") + 1 }

func isGlob(p string) bool { return strings.ContainsAny(p, "*?[") }

// Scan walks the project root and returns all matching signatures.
func Scan(root string) ([]Detection, error) {
	entries := map[string]bool{} // rel slash path -> isDir

	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			// Ignore permission errors on individual entries.
			return nil
		}
		rel, rerr := filepath.Rel(root, p)
		if rerr != nil {
			return nil
		}
		if rel == "." {
			return nil
		}
		slash := filepath.ToSlash(rel)
		entries[slash] = d.IsDir()

		if d.IsDir() {
			base := filepath.Base(p)
			if heavyDirs[base] {
				return fs.SkipDir
			}
			if depth(slash) > 3 {
				return fs.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var out []Detection
	for _, sig := range Signatures {
		var matched []string
		for _, m := range sig.Markers {
			if markerPresent(m, entries) {
				matched = append(matched, m.Path)
			}
		}
		if len(matched) > 0 {
			out = append(out, Detection{Signature: sig, Matched: matched})
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		return out[i].Key < out[j].Key
	})
	return out, nil
}

func markerPresent(m Marker, entries map[string]bool) bool {
	if isGlob(m.Path) {
		for rel, isDir := range entries {
			base := path.Base(rel)
			if ok, _ := path.Match(m.Path, base); ok {
				if !m.Dir || isDir {
					return true
				}
			}
			// Also allow matching the full relative path (e.g. "*.sln" at root).
			if ok, _ := path.Match(m.Path, rel); ok {
				if !m.Dir || isDir {
					return true
				}
			}
		}
		return false
	}
	isDir, ok := entries[m.Path]
	if !ok {
		return false
	}
	return !m.Dir || isDir
}

// Keys returns the ordered, unique template keys from detections.
func Keys(ds []Detection) []string {
	seen := map[string]bool{}
	var keys []string
	for _, d := range ds {
		if !seen[d.Key] {
			seen[d.Key] = true
			keys = append(keys, d.Key)
		}
	}
	return keys
}
