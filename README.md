<div align="center">

# 🦞 ignorekit

**Offline, zero-dependency `.gitignore` generator & Git tracking-hygiene workbench**

Detect your stack → generate a correct `.gitignore` → find files you already
committed but should ignore → fix it in one command.

🌐 **Language**
[English](README.md) · [简体中文](docs/README.zh-CN.md) · [繁體中文](docs/README.zh-TW.md)

⬇️ **[Download the latest release](https://github.com/gitstq/ignorekit/releases/latest)**

</div>

---

## 🎉 Introduction

AI coding assistants and agents are everywhere now — **Claude Code, Cursor,
Windsurf, GitHub Copilot, Aider, Continue, Cline, Roo, Codex, MCP servers** and
many more. They leave a trail of files inside your repository: session history,
agent memory, local settings, caches, and MCP configs that can contain **API
keys and credentials**.

A single careless `git add .` can leak that data. At the same time, *some* of
these files are meant to be shared with your team (for example `.cursorrules`,
`CLAUDE.md`, or `copilot-instructions.md`). Telling the two apart is tedious and
error-prone — and classic `.gitignore` generators are slow to cover new AI tools
and never tell you that **you already committed something you shouldn't have**.

**ignorekit** solves this. It is a fast, single-binary, fully offline tool that:

- 🧠 understands modern languages, frameworks, editors and **AI agents / MCP**,
  distinguishing files you should *commit* from files you should *ignore*;
- 🧪 cross-checks the suggested rules against Git's tracked files, surfacing
  **tracked-but-ignored** files and high-confidence secret red flags;
- 🩹 merges rules safely into your existing `.gitignore` and can untrack files
  with `git rm --cached`, idempotently and without touching your working tree.

> 🍪 Inspiration: the current wave of agent tooling on GitHub Trending
> (agent managers, agent memory, MCP servers) highlighted how much local state
> these tools create. ignorekit is independently designed and implemented; it
> copies no existing project, and ignore patterns are universal facts.

---

## ✨ Features

- 🔌 **Zero runtime dependencies** — a single static binary, written in pure Go
  (standard library only). No Node, Python, Docker or network required.
- 📴 **Fully offline** — every template is embedded. Nothing is ever uploaded.
- 🤖 **AI agent & MCP aware** — Claude Code, Cursor, Windsurf, Codeium, Copilot,
  Aider, Continue, Cline, Roo, Cody, Codex, Gemini CLI, OpenCode, Tabnine and
  local MCP configs, with explicit **"keep/commit" vs "ignore"** guidance.
- 🧰 **Broad stack coverage** — 60+ templates across operating systems, editors
  (VS Code, JetBrains, Vim, Emacs, Eclipse, Sublime), languages (Node, Python,
  Go, Rust, Java, .NET, Ruby, PHP, Dart/Flutter, Swift, Elixir, Haskell, R),
  frameworks (Next, Nuxt, SvelteKit, Astro, Vite, Angular, Gatsby, Parcel,
  Turborepo, Storybook, Playwright/Cypress, Expo, Unity, Unreal) and more.
- 🧹 **Tracking-hygiene audit** — finds files Git already tracks that match an
  ignore rule, the capability classic generators lack.
- 🔐 **Secret red flags** — a small, high-confidence set of credential patterns
  (GitHub/GitLab/AWS/Google/Slack/Stripe/OpenAI tokens, JWTs, private keys).
  It is deliberately not a full secret scanner — pair it with dedicated tools.
- 🧷 **Safe, idempotent merge** — preserves your hand-written rules, deduplicates,
  groups by section, and tracks its own rules inside a clearly marked block.
- 🛠️ **Fix mode** — write the file and optionally `git rm --cached` in one step;
  working-tree files are always kept.
- 🚦 **CI friendly** — `text`, `json` and `SARIF 2.1.0` output plus stable exit
  codes (0 clean, 1 findings, 2 error).
- 🪟 **Cross-platform** — Windows, macOS and Linux, `amd64`/`arm64` (and `386`).

---

## 🚀 Quick Start

### Requirements

None beyond the binary. Prebuilt executables are available for Windows, macOS
and Linux. To build from source you need **Go 1.23+**.

### Option A — Download a prebuilt binary (recommended)

1. Open the **[latest release](https://github.com/gitstq/ignorekit/releases/latest)**.
2. Download the archive for your OS/architecture.
3. Verify the checksum (optional) and put `ignorekit` (or `ignorekit.exe`) on your `PATH`.

```bash
# Example for Linux x86_64
tar -xzf ignorekit_1.0.0_linux_amd64.tar.gz
sudo mv ignorekit_1.0.0_linux_amd64/ignorekit /usr/local/bin/
ignorekit version
```

### Option B — Install with Go

```bash
go install github.com/gitstq/ignorekit@latest
```

### Option C — Build from source

```bash
git clone https://github.com/gitstq/ignorekit.git
cd ignorekit
make build          # or: go build -trimpath -o ignorekit .
```

### 30-second tour

```bash
cd your-project

# 1) See the .gitignore ignorekit would generate (prints to stdout)
ignorekit init

# 2) Read-only hygiene audit: tracked files to ignore + secret red flags
ignorekit audit

# 3) Apply: merge rules into .gitignore and untrack files that should be ignored
ignorekit apply --untrack
```

---

## 📖 Usage Guide

### Commands

| Command | Description |
|---|---|
| `init` | Detect tools and generate a `.gitignore` (stdout by default) |
| `audit` | Read-only hygiene report (tracked-but-ignored files + secrets) |
| `apply` | Merge suggested rules into `.gitignore` (safe & idempotent) |
| `list-templates` | List every built-in template |
| `version` | Print version information |

### Global flags

| Flag | Description |
|---|---|
| `-C`, `--dir DIR` | Project directory (default: current directory) |
| `-f`, `--format FMT` | Output format: `text` \| `json` \| `sarif` (audit) |
| `--write` | Write the result to `.gitignore` (init) |
| `--force` | Overwrite an existing `.gitignore` (init) |
| `--untrack` | Also run `git rm --cached` on tracked-but-ignored files (apply) |
| `--dry-run` | Show the plan without writing (apply) |
| `-y`, `--yes` | Proceed without prompting |
| `--no-color` | Disable colored output |
| `--strict` | Treat a missing `.gitignore` with suggestions as a finding (audit) |

Global flags work **before or after** the command, e.g. both
`ignorekit -C ./app audit` and `ignorekit audit -C ./app` are valid.

### `init` — generate a `.gitignore`

```bash
ignorekit init                 # preview on stdout
ignorekit init --write         # create .gitignore (fails if it already exists)
ignorekit init --write --force # overwrite an existing file
```

### `audit` — read-only hygiene report

```bash
ignorekit audit                       # human-readable
ignorekit audit --format json         # machine-readable
ignorekit audit --format sarif        # SARIF 2.1.0 for code scanning
```

The audit reports, per detected tool, which files you should **keep/commit**
(green) and which need **review** (yellow). It then lists:

- tracked files that match an ignore rule (new vs existing rules);
- high-confidence secret red flags with file and line number (redacted).

### `apply` — merge and optionally untrack

```bash
ignorekit apply                 # merge suggested rules (creates a managed block)
ignorekit apply --dry-run       # preview what would change
ignorekit apply --untrack       # also remove tracked-but-ignored files from the index
```

ignorekit writes its rules inside a clearly delimited block:

```gitignore
# >>> ignorekit managed (safe to delete; regenerate with `ignorekit apply`) >>>
# --- Node.js ---
node_modules/
# <<< ignorekit managed <<<
```

Your hand-written rules above the block are never modified or duplicated, and
re-running `apply` is safe: it updates the block in place and reports
`0 newly added` when nothing changed.

### AI tool guidance (what to commit vs ignore)

| Tool | Commit / share | Ignore |
|---|---|---|
| Claude Code | `CLAUDE.md`, `.claude/settings.json` | `.claude/settings.local.json`, todos, shell snapshots, session/project state |
| Cursor | `.cursorrules`, `.cursor/rules/` | cache, index, eval, sessions, `*.local` |
| Windsurf | `.windsurf/rules/` | cache, index, sessions, `*.local` |
| GitHub Copilot | `.github/copilot-instructions.md` | local Copilot state |
| Aider | — | input/LLM/chat history, `.aider.conf.yml`, caches |
| Continue / Cline / Roo / Cody / Codex | rule files (e.g. `.clinerules`, `.roorules`, `AGENTS.md`) | local state directories |
| MCP | shared server config (after review) | `*.mcp.local`, `mcp.local.json` |

> Shared MCP config files (`.mcp.json`, `mcp.json`) are **not** auto-ignored
> because teams often version them. ignorekit instead scans them for secrets and
> flags them for review. Use local variants (`mcp.local.json`) for machine-only
> credentials.

### CI integration

`audit` exits `1` when tracked-hygiene or secret findings exist, `0` when clean,
and `2` on error — perfect for a CI gate.

**GitHub Actions** (upload SARIF to code scanning):

```yaml
name: ignorekit
on: [push, pull_request]
jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run ignorekit
        run: |
          curl -L -o ignorekit.tgz \
            https://github.com/gitstq/ignorekit/releases/latest/download/ignorekit_1.0.0_linux_amd64.tar.gz
          tar -xzf ignorekit.tgz
          ./ignorekit_1.0.0_linux_amd64/ignorekit audit --format sarif > ignorekit.sarif
      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: ignorekit.sarif
```

### Typical scenarios

- 🆕 **New project** — `ignorekit init --write` for a complete baseline.
- 🧑‍🤝‍🧑 **Onboarding an existing repo** — `ignorekit audit`, then
  `ignorekit apply --untrack`, review the staged deletions, and commit.
- 🔒 **Pre-push safety net** — run `ignorekit audit` before pushing or in CI.
- 🤖 **After installing an AI assistant** — re-run `apply` to cover its local state.

### Screenshots / demo

> 📸 Placeholder — add a terminal recording here, e.g.
> `![demo](docs/demo.gif)`.

---

## 💡 Design & Roadmap

### Design principles

- **Offline & local-first.** Secrets and source code should never leave your
  machine; templates are embedded and matching is implemented in pure Go.
- **Zero dependencies.** A single static binary removes installation friction
  and supply-chain risk.
- **Safe by default.** All mutating commands support a dry-run, never delete
  working-tree files, and isolate generated rules in a removable block.
- **High signal over noise.** Secret detection favors confidence to avoid
  alert fatigue; we defer comprehensive scanning to specialized tools.

### Why Go?

Go produces small, static, cross-platform binaries from one codebase, has a rich
standard library (including regexp and flag parsing), and makes the
cross-compilation of all release targets trivial with `CGO_ENABLED=0`.

### Roadmap

- [ ] More AI tools and framework signatures (community-driven).
- [ ] Custom user templates via an optional config file.
- [ ] An interactive `--interactive` mode for reviewing each finding.
- [ ] Pre-commit hook installer and more CI templates.
- [ ] Optional deeper secret/entropy checks behind an explicit flag.
- [ ] Japanese, Korean and Spanish README translations.

Contributions and roadmap votes via issues and PRs are very welcome — see
[CONTRIBUTING.md](CONTRIBUTING.md).

---

## 📦 Packaging & Deployment

ignorekit is a **desktop/CLI application** shipped as prebuilt binaries.

### Build all release archives

```bash
make cross          # or: VERSION=1.0.0 ./scripts/build.sh
```

This cross-compiles every target with `CGO_ENABLED=0`, `-trimpath` and stripped
symbols, then writes archives to `dist/`:

- Windows: `.zip` (`amd64`, `arm64`, `386`)
- macOS / Linux: `.tar.gz` (`amd64`, `arm64`; Linux also `386`)
- `dist/SHA256SUMS` with SHA-256 checksums of every archive

### Verify a download

```bash
sha256sum -c SHA256SUMS --ignore-missing
```

### Compatibility

| OS | Architectures |
|---|---|
| Windows | amd64, arm64, 386 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Linux | amd64, arm64, 386 |

The binaries are statically linked and have no runtime library dependencies.

---

## 🤝 Contributing

Pull requests and issues are welcome! Please read
[CONTRIBUTING.md](CONTRIBUTING.md) for the full workflow. In short:

1. Fork and branch from `main`.
2. Add tests and ensure `go test ./...` passes; run `gofmt -s` and `go vet`.
3. Use [Conventional Commits](https://www.conventionalcommits.org/)
   (`feat:`, `fix:`, `docs:`, `refactor:`).
4. Keep the project dependency-free unless discussed in an issue.

Please be kind and constructive in all discussions.

---

## 📄 License

ignorekit is released under the **[MIT License](LICENSE)**.

<div align="center">
🦞 Built for developers who'd rather not leak their agent's secrets.
</div>
