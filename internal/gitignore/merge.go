package gitignore

import (
	"strings"
)

// Markers delimit the block of rules owned by ignorekit so repeated runs
// are idempotent and never duplicate rules.
const (
	BeginMarker = "# >>> ignorekit managed (safe to delete; regenerate with `ignorekit apply`) >>>"
	EndMarker   = "# <<< ignorekit managed <<<"
)

// Rule is a single generated rule line bound to a section.
type Rule struct {
	Line    string
	Section string
}

// Signature returns a normalized identity for de-duplication.
func (p *Pattern) Signature() string { return p.sig }

// ruleSignature compiles a single line and returns its signature.
func ruleSignature(line string) string {
	p, err := Compile(line, "")
	if err != nil || p == nil {
		return strings.TrimSpace(line)
	}
	return p.sig
}

// SplitManaged separates existing content into the user-owned prefix and the
// body of a previously written managed block (if any).
func SplitManaged(existing string) (userContent, managedBody string, found bool) {
	lines := strings.Split(existing, "\n")
	var user, managed []string
	inBlock := false
	for _, ln := range lines {
		trim := strings.TrimSpace(ln)
		switch {
		case trim == BeginMarker:
			inBlock = true
		case trim == EndMarker:
			inBlock = false
			found = true
		case inBlock:
			managed = append(managed, ln)
		default:
			user = append(user, ln)
		}
	}
	u := strings.TrimRight(strings.Join(user, "\n"), "\n")
	m := strings.Join(managed, "\n")
	return u, m, found
}

// existingSignatures collects normalized rule signatures already present in
// user-owned content and (optionally) the prior managed block.
func existingSignatures(userContent, managedBody string) map[string]bool {
	sigs := map[string]bool{}
	for _, src := range []string{userContent, managedBody} {
		fs := NewFileSet()
		_ = fs.ParseReader(src, "")
		for _, p := range fs.Patterns {
			sigs[p.sig] = true
		}
	}
	return sigs
}

// MergeResult describes the outcome of a merge.
type MergeResult struct {
	Content     string   // full merged file content
	Added       []string // rules newly introduced (not present anywhere before)
	InManaged   []string // rules present in the managed block after merge
	HadManaged  bool
	UserContent string
}

// RenderManaged builds the managed block from the desired rules, omitting any
// rule already declared in user-owned content, and grouping by section.
func RenderManaged(desired []Rule, userSigs map[string]bool) (blockBody string, included []string) {
	var b strings.Builder
	lastSection := ""
	for _, r := range desired {
		sig := ruleSignature(r.Line)
		if userSigs[sig] {
			continue // already hand-declared outside the managed block
		}
		if r.Section != lastSection {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString("# --- " + r.Section + " ---\n")
			lastSection = r.Section
		}
		b.WriteString(r.Line + "\n")
		included = append(included, r.Line)
	}
	return strings.TrimRight(b.String(), "\n"), included
}

// Merge combines desired generated rules with an existing .gitignore.
func Merge(existing string, desired []Rule) MergeResult {
	userContent, oldManaged, had := SplitManaged(existing)
	userSigs := existingSignatures(userContent, "")
	allSigs := existingSignatures(userContent, oldManaged)

	body, included := RenderManaged(desired, userSigs)

	var added []string
	for _, r := range desired {
		if !allSigs[ruleSignature(r.Line)] {
			added = append(added, r.Line)
		}
	}

	var b strings.Builder
	u := strings.TrimSpace(userContent)
	if u != "" {
		b.WriteString(u)
		b.WriteString("\n\n")
	}
	if strings.TrimSpace(body) != "" {
		b.WriteString(BeginMarker + "\n")
		b.WriteString(body + "\n")
		b.WriteString(EndMarker + "\n")
	}
	res := MergeResult{
		Content:     b.String(),
		Added:       added,
		InManaged:   included,
		HadManaged:  had,
		UserContent: userContent,
	}
	return res
}
