// Package detect scans a project directory and identifies languages,
// frameworks, editors, AI coding assistants/agents, MCP setups and other
// tooling by the files and directories they leave behind.
package detect

// Marker is a file or directory whose presence indicates a tool.
type Marker struct {
	Path string // slash-separated, relative to the project root
	Dir  bool
}

// Signature describes one detectable tool.
type Signature struct {
	Key      string // matches a templates.Template key
	Name     string
	Category string
	Markers  []Marker
	// Shared lists files that should normally be committed (guidance only).
	Shared []string
	// Review lists files that may contain secrets and need manual review.
	Review []string
	Notes  string
}

// Signatures is the full detection table.
var Signatures = []Signature{
	// Languages / runtimes
	{Key: "node", Name: "Node.js", Category: "Language", Markers: []Marker{
		{Path: "package.json"}, {Path: "package-lock.json"}, {Path: "yarn.lock"},
		{Path: "pnpm-lock.yaml"}, {Path: "bun.lockb"}, {Path: "node_modules", Dir: true},
	}},
	{Key: "python", Name: "Python", Category: "Language", Markers: []Marker{
		{Path: "requirements.txt"}, {Path: "pyproject.toml"}, {Path: "setup.py"},
		{Path: "setup.cfg"}, {Path: "Pipfile"}, {Path: "poetry.lock"},
		{Path: "environment.yml"}, {Path: "tox.ini"}, {Path: "__pycache__", Dir: true},
	}},
	{Key: "go", Name: "Go", Category: "Language", Markers: []Marker{
		{Path: "go.mod"}, {Path: "go.sum"},
	}},
	{Key: "rust", Name: "Rust", Category: "Language", Markers: []Marker{
		{Path: "Cargo.toml"}, {Path: "Cargo.lock"}, {Path: "target", Dir: true},
	}},
	{Key: "java", Name: "Java / Gradle / Maven", Category: "Language", Markers: []Marker{
		{Path: "pom.xml"}, {Path: "build.gradle"}, {Path: "build.gradle.kts"},
		{Path: "settings.gradle"}, {Path: "target", Dir: true}, {Path: "gradlew"},
	}},
	{Key: "dotnet", Name: ".NET / C#", Category: "Language", Markers: []Marker{
		{Path: "global.json"}, {Path: "*.sln"}, {Path: "*.csproj"},
	}},
	{Key: "ruby", Name: "Ruby", Category: "Language", Markers: []Marker{
		{Path: "Gemfile"}, {Path: "*.gemspec"}, {Path: ".ruby-version"},
	}},
	{Key: "php", Name: "PHP / Composer", Category: "Language", Markers: []Marker{
		{Path: "composer.json"}, {Path: "composer.lock"},
	}},
	{Key: "dart", Name: "Dart / Flutter", Category: "Language", Markers: []Marker{
		{Path: "pubspec.yaml"}, {Path: "pubspec.lock"}, {Path: ".dart_tool", Dir: true},
	}},
	{Key: "swift", Name: "Swift / Xcode", Category: "Language", Markers: []Marker{
		{Path: "Package.swift"}, {Path: "*.xcodeproj", Dir: true}, {Path: "*.xcworkspace", Dir: true},
	}},
	{Key: "elixir", Name: "Elixir", Category: "Language", Markers: []Marker{
		{Path: "mix.exs"}, {Path: "mix.lock"},
	}},
	{Key: "haskell", Name: "Haskell", Category: "Language", Markers: []Marker{
		{Path: "stack.yaml"}, {Path: "*.cabal"}, {Path: "cabal.project"},
	}},
	{Key: "r-lang", Name: "R", Category: "Language", Markers: []Marker{
		{Path: "*.R"}, {Path: "*.Rmd"}, {Path: "DESCRIPTION"},
	}},

	// Frameworks
	{Key: "next", Name: "Next.js", Category: "Framework", Markers: []Marker{
		{Path: ".next", Dir: true}, {Path: "next.config.js"}, {Path: "next.config.mjs"}, {Path: "next.config.ts"},
	}},
	{Key: "nuxt", Name: "Nuxt", Category: "Framework", Markers: []Marker{
		{Path: ".nuxt", Dir: true}, {Path: ".output", Dir: true}, {Path: "nuxt.config.ts"}, {Path: "nuxt.config.js"},
	}},
	{Key: "svelte", Name: "SvelteKit", Category: "Framework", Markers: []Marker{
		{Path: ".svelte-kit", Dir: true}, {Path: "svelte.config.js"}, {Path: "svelte.config.ts"},
	}},
	{Key: "astro", Name: "Astro", Category: "Framework", Markers: []Marker{
		{Path: ".astro", Dir: true}, {Path: "astro.config.mjs"}, {Path: "astro.config.ts"},
	}},
	{Key: "vite", Name: "Vite", Category: "Framework", Markers: []Marker{
		{Path: "vite.config.ts"}, {Path: "vite.config.js"}, {Path: "node_modules/.vite", Dir: true},
	}},
	{Key: "angular", Name: "Angular", Category: "Framework", Markers: []Marker{
		{Path: ".angular", Dir: true}, {Path: "angular.json"},
	}},
	{Key: "gatsby", Name: "Gatsby", Category: "Framework", Markers: []Marker{
		{Path: ".cache", Dir: true}, {Path: "gatsby-config.js"}, {Path: "gatsby-config.ts"},
	}},
	{Key: "parcel", Name: "Parcel", Category: "Framework", Markers: []Marker{
		{Path: ".parcel-cache", Dir: true},
	}},
	{Key: "turbo", Name: "Turborepo", Category: "Framework", Markers: []Marker{
		{Path: "turbo.json"}, {Path: ".turbo", Dir: true},
	}},
	{Key: "tsc", Name: "TypeScript / ESLint cache", Category: "Framework", Markers: []Marker{
		{Path: "tsconfig.json"}, {Path: ".eslintcache"},
	}},
	{Key: "storybook", Name: "Storybook", Category: "Framework", Markers: []Marker{
		{Path: ".storybook", Dir: true}, {Path: "storybook-static", Dir: true},
	}},
	{Key: "test-e2e", Name: "Playwright / Cypress", Category: "Framework", Markers: []Marker{
		{Path: "playwright.config.ts"}, {Path: "cypress.config.ts"}, {Path: "test-results", Dir: true},
		{Path: "playwright-report", Dir: true}, {Path: "cypress", Dir: true},
	}},
	{Key: "expo", Name: "Expo", Category: "Framework", Markers: []Marker{
		{Path: ".expo", Dir: true}, {Path: "app.json"}, {Path: "expo.json"},
	}},
	{Key: "jekyll", Name: "Jekyll / Hexo", Category: "Framework", Markers: []Marker{
		{Path: "_config.yml"}, {Path: "_config.toml"}, {Path: "_site", Dir: true},
	}},
	{Key: "unity", Name: "Unity", Category: "Framework", Markers: []Marker{
		{Path: "Library", Dir: true}, {Path: "Assets", Dir: true}, {Path: "ProjectSettings", Dir: true},
	}},
	{Key: "unreal", Name: "Unreal Engine", Category: "Framework", Markers: []Marker{
		{Path: "Content", Dir: true}, {Path: "Source", Dir: true}, {Path: "*.uproject"},
	}},

	// Editors / IDEs
	{Key: "vscode", Name: "Visual Studio Code", Category: "Editor", Markers: []Marker{
		{Path: ".vscode", Dir: true},
	}, Shared: []string{
		".vscode/settings.json", ".vscode/tasks.json", ".vscode/launch.json", ".vscode/extensions.json",
	}},
	{Key: "jetbrains", Name: "JetBrains IDEs", Category: "Editor", Markers: []Marker{
		{Path: ".idea", Dir: true}, {Path: "*.iml"},
	}},
	{Key: "vim", Name: "Vim", Category: "Editor", Markers: []Marker{
		{Path: ".vimrc"}, {Path: "*.swp"},
	}},
	{Key: "emacs", Name: "Emacs", Category: "Editor", Markers: []Marker{
		{Path: ".emacs"}, {Path: "init.el"},
	}},
	{Key: "eclipse", Name: "Eclipse / STS", Category: "Editor", Markers: []Marker{
		{Path: ".project"}, {Path: ".classpath"}, {Path: ".settings", Dir: true},
	}},
	{Key: "sublime", Name: "Sublime Text", Category: "Editor", Markers: []Marker{
		{Path: "*.sublime-project"}, {Path: "*.sublime-workspace"},
	}},
	{Key: "scala-editors", Name: "Metals / Bloop", Category: "Editor", Markers: []Marker{
		{Path: ".metals", Dir: true}, {Path: ".bloop", Dir: true}, {Path: ".bsp", Dir: true},
	}},

	// AI assistants / agents / MCP
	{Key: "claude-code", Name: "Claude Code", Category: "AI", Markers: []Marker{
		{Path: ".claude", Dir: true}, {Path: ".claude.json"}, {Path: "CLAUDE.md"},
	}, Shared: []string{"CLAUDE.md", ".claude/settings.json"}, Notes: "CLAUDE.md 与共享 settings.json 应提交；本地设置与会话历史应忽略。"},
	{Key: "cursor", Name: "Cursor", Category: "AI", Markers: []Marker{
		{Path: ".cursor", Dir: true}, {Path: ".cursorrules"},
	}, Shared: []string{".cursorrules", ".cursor/rules/"}, Notes: ".cursorrules 与 .cursor/rules 规则应共享；缓存/索引/会话应忽略。"},
	{Key: "windsurf", Name: "Windsurf", Category: "AI", Markers: []Marker{
		{Path: ".windsurf", Dir: true},
	}, Shared: []string{".windsurf/rules/"}},
	{Key: "codeium", Name: "Codeium", Category: "AI", Markers: []Marker{
		{Path: ".codeium", Dir: true},
	}},
	{Key: "copilot", Name: "GitHub Copilot", Category: "AI", Markers: []Marker{
		{Path: ".github/copilot-instructions.md"}, {Path: ".copilot", Dir: true},
	}, Shared: []string{".github/copilot-instructions.md"}},
	{Key: "aider", Name: "Aider", Category: "AI", Markers: []Marker{
		{Path: ".aider.conf.yml"}, {Path: ".aider.input.history"}, {Path: ".aider.llm.history"},
		{Path: ".aider.chat.history.md"},
	}},
	{Key: "continue", Name: "Continue", Category: "AI", Markers: []Marker{
		{Path: ".continue", Dir: true},
	}},
	{Key: "cline", Name: "Cline", Category: "AI", Markers: []Marker{
		{Path: ".cline", Dir: true}, {Path: ".clinerules"},
	}, Shared: []string{".clinerules"}},
	{Key: "roo", Name: "Roo Code", Category: "AI", Markers: []Marker{
		{Path: ".roo", Dir: true}, {Path: ".roorules"},
	}, Shared: []string{".roorules"}},
	{Key: "cody", Name: "Sourcegraph Cody", Category: "AI", Markers: []Marker{
		{Path: ".cody", Dir: true},
	}},
	{Key: "codex", Name: "OpenAI Codex CLI", Category: "AI", Markers: []Marker{
		{Path: ".codex", Dir: true}, {Path: "AGENTS.md"},
	}, Shared: []string{"AGENTS.md"}},
	{Key: "gemini-cli", Name: "Gemini CLI", Category: "AI", Markers: []Marker{
		{Path: ".gemini", Dir: true}, {Path: "GEMINI.md"},
	}, Shared: []string{"GEMINI.md"}},
	{Key: "opencode", Name: "OpenCode", Category: "AI", Markers: []Marker{
		{Path: ".opencode", Dir: true}, {Path: "opencode.json"},
	}, Shared: []string{"opencode.json"}},
	{Key: "v0", Name: "v0", Category: "AI", Markers: []Marker{
		{Path: ".v0", Dir: true},
	}},
	{Key: "mcp-local", Name: "MCP local configs", Category: "AI", Markers: []Marker{
		{Path: "mcp.local.json"}, {Path: ".mcp.local.json"},
	}},
	{Key: "tabnine", Name: "Tabnine", Category: "AI", Markers: []Marker{
		{Path: ".tabnine", Dir: true},
	}},

	// Environment & secrets
	{Key: "env", Name: "Environment variables", Category: "Secrets", Markers: []Marker{
		{Path: ".env"}, {Path: ".env.local"}, {Path: ".env.development"}, {Path: ".env.production"},
	}, Review: []string{".env"}, Shared: []string{".env.example", ".env.template", ".env.sample"}},
	{Key: "cloud-keys", Name: "Cloud credentials", Category: "Secrets", Markers: []Marker{
		{Path: ".aws", Dir: true}, {Path: ".gcp", Dir: true}, {Path: ".azure", Dir: true},
		{Path: "*.pem"}, {Path: "*.key"},
	}},
	{Key: "kubeconfig", Name: "Kubeconfig", Category: "Secrets", Markers: []Marker{
		{Path: "kubeconfig"}, {Path: ".kube", Dir: true},
	}},
	{Key: "terraform", Name: "Terraform", Category: "Secrets", Markers: []Marker{
		{Path: "*.tf"}, {Path: ".terraform", Dir: true}, {Path: "*.tfstate"},
	}},
	{Key: "ansible", Name: "Ansible", Category: "Secrets", Markers: []Marker{
		{Path: "ansible.cfg"}, {Path: "group_vars", Dir: true}, {Path: "host_vars", Dir: true},
		{Path: "playbooks", Dir: true},
	}},
	{Key: "docker", Name: "Docker", Category: "Deploy", Markers: []Marker{
		{Path: "Dockerfile"}, {Path: "docker-compose.yml"}, {Path: "docker-compose.yaml"},
	}},
	{Key: "vercel", Name: "Vercel", Category: "Deploy", Markers: []Marker{
		{Path: ".vercel", Dir: true}, {Path: "vercel.json"},
	}},
	{Key: "netlify", Name: "Netlify", Category: "Deploy", Markers: []Marker{
		{Path: ".netlify", Dir: true}, {Path: "netlify.toml"},
	}},
	{Key: "serverless", Name: "Serverless / SST / Wrangler", Category: "Deploy", Markers: []Marker{
		{Path: ".serverless", Dir: true}, {Path: ".sst", Dir: true}, {Path: ".wrangler", Dir: true},
		{Path: "wrangler.toml"},
	}},
	{Key: "firebase", Name: "Firebase", Category: "Deploy", Markers: []Marker{
		{Path: ".firebase", Dir: true}, {Path: "firebase.json"},
	}},
	{Key: "logs", Name: "Logs", Category: "Logs", Markers: []Marker{
		{Path: "*.log"},
	}},
	{Key: "temp", Name: "Temporary / cache", Category: "Logs", Markers: []Marker{
		{Path: ".tmp", Dir: true}, {Path: "tmp", Dir: true},
	}},
}
