// Package templates holds the built-in, offline .gitignore rule templates.
// Everything is compiled into the binary; no network access is ever required.
package templates

// Template is a named, composable set of ignore rules.
type Template struct {
	Key      string
	Title    string
	Category string
	Order    int
	Rules    []string
}

// Category ordering for stable output.
const (
	CatOS       = "Operating System"
	CatEditor   = "Editor / IDE"
	CatLanguage = "Language / Runtime"
	CatFramework = "Framework / Build"
	CatAI       = "AI Assistant / Agent / MCP"
	CatSecrets  = "Environment & Secrets"
	CatDeploy   = "Deployment / Cloud"
	CatLogs     = "Logs & Temp"
)

var registry = map[string]Template{
	// ---------------- Operating systems ----------------
	"osx": {Key: "osx", Title: "macOS", Category: CatOS, Order: 10, Rules: []string{
		".DS_Store", ".AppleDouble", ".LSOverride", "Icon?", "._*",
		".Spotlight-V100", ".Trashes", ".fseventsd", ".TemporaryItems",
		".VolumeIcon.icns", ".com.apple.timemachine.donotpresent", ".AppleDB",
		".AppleDesktop", ".apdisk",
	}},
	"windows": {Key: "windows", Title: "Windows", Category: CatOS, Order: 11, Rules: []string{
		"Thumbs.db", "Thumbs.db:encryptable", "ehthumbs.db", "ehthumbs_vista.db",
		"Desktop.ini", "$RECYCLE.BIN/", "*.cab", "*.msi", "*.msix", "*.msm", "*.msp",
		"*.lnk",
	}},
	"linux": {Key: "linux", Title: "Linux", Category: CatOS, Order: 12, Rules: []string{
		"*~", ".fuse_hidden*", ".directory", ".Trash-*", ".nfs*",
	}},

	// ---------------- Editors / IDEs ----------------
	"vscode": {Key: "vscode", Title: "Visual Studio Code", Category: CatEditor, Order: 20, Rules: []string{
		".vscode/*",
		"!.vscode/settings.json",
		"!.vscode/tasks.json",
		"!.vscode/launch.json",
		"!.vscode/extensions.json",
		"!.vscode/*.code-snippets",
		".vscode/*.local",
		".history/",
	}},
	"jetbrains": {Key: "jetbrains", Title: "JetBrains IDEs", Category: CatEditor, Order: 21, Rules: []string{
		".idea/", "*.iml", "*.iws", "*.ipr", "*.idea_modules/", "out/",
	}},
	"vim": {Key: "vim", Title: "Vim", Category: CatEditor, Order: 22, Rules: []string{
		"*.swp", "*.swo", "*.un~", "[._]*.s[a-v][a-z]", "Session.vim", "tags",
	}},
	"emacs": {Key: "emacs", Title: "Emacs", Category: CatEditor, Order: 23, Rules: []string{
		"*~", ".\\#*", ".#*", "*.elc", "auto-save-list/", "tramp", ".emacs.desktop",
		".emacs.desktop.lock",
	}},
	"eclipse": {Key: "eclipse", Title: "Eclipse / STS", Category: CatEditor, Order: 24, Rules: []string{
		".metadata/", ".tmp/", ".bak/", "*.tmp_user", ".localhistory/",
		".project", ".classpath", ".settings/", ".factorypath", ".springBeans",
	}},
	"sublime": {Key: "sublime", Title: "Sublime Text", Category: CatEditor, Order: 25, Rules: []string{
		"*.sublime-workspace", "*.sublime-project", "sftp-config.json",
	}},
	"scala-editors": {Key: "scala-editors", Title: "Scala Metals / Bloop", Category: CatEditor, Order: 26, Rules: []string{
		".bloop/", ".metals/", ".bsp/", ".scala_server/",
	}},

	// ---------------- Languages / runtimes ----------------
	"node": {Key: "node", Title: "Node.js", Category: CatLanguage, Order: 30, Rules: []string{
		"node_modules/", "npm-debug.log*", "yarn-debug.log*", "yarn-error.log*",
		".npm/", ".pnpm-debug.log*", "lerna-debug.log*", ".node_repl_history",
	}},
	"python": {Key: "python", Title: "Python", Category: CatLanguage, Order: 31, Rules: []string{
		"__pycache__/", "*.py[cod]", "*$py.class", "*.so", ".Python",
		"venv/", ".venv/", "env/", "ENV/", ".eggs/", "*.egg", "*.egg-info/",
		".pytest_cache/", ".mypy_cache/", ".ruff_cache/", ".pyre/", ".pytype/",
		".coverage", ".coverage.*", "htmlcov/", ".tox/", ".nox/", ".hypothesis/",
		"coverage.xml", "*.cover", ".ipynb_checkpoints/", "profile_default/",
		"instance/", ".webassets-cache/",
	}},
	"go": {Key: "go", Title: "Go", Category: CatLanguage, Order: 32, Rules: []string{
		"*.exe", "*.exe~", "*.dll", "*.so", "*.dylib", "*.test", "*.out", "vendor/",
	}},
	"rust": {Key: "rust", Title: "Rust", Category: CatLanguage, Order: 33, Rules: []string{
		"target/", "**/*.rs.bk", "*.pdb", "Cargo.lock.bak",
	}},
	"java": {Key: "java", Title: "Java / Gradle / Maven", Category: CatLanguage, Order: 34, Rules: []string{
		"target/", "*.class", "*.war", "*.nar", "*.ear", ".gradle/", "build/",
		"hs_err_pid*", "replay_pid*", ".mvn/timing.properties", ".mvn/wrapper/maven-wrapper.jar",
		"dependency-reduced-pom.xml",
	}},
	"dotnet": {Key: "dotnet", Title: ".NET / C#", Category: CatLanguage, Order: 35, Rules: []string{
		"[Bb]in/", "[Oo]bj/", "[Dd]ebug/", "[Rr]elease/", "*.user", "*.suo",
		"*.userosscache", "*.sln.docstates", "project.lock.json", "project.fragment.lock",
		"project.assets.json", "*.nuget.props", "*.nuget.targets", ".vs/",
	}},
	"ruby": {Key: "ruby", Title: "Ruby / Bundler", Category: CatLanguage, Order: 36, Rules: []string{
		"/.bundle/", "/vendor/bundle/", "*.gem", ".gem/", "log/", "tmp/",
		".byebug_history", ".ruby-gemset", ".bundle",
	}},
	"php": {Key: "php", Title: "PHP / Composer", Category: CatLanguage, Order: 37, Rules: []string{
		"vendor/", "composer.phar", ".phpunit.result.cache", ".phpunit.cache/",
		"Homestead.json", "Homestead.yaml", ".phpactor.json", ".php-cs-fixer.cache",
		"node_modules/",
	}},
	"dart": {Key: "dart", Title: "Dart / Flutter", Category: CatLanguage, Order: 38, Rules: []string{
		".dart_tool/", ".flutter-plugins", ".flutter-plugins-dependencies", ".packages",
		".pub-cache/", ".pub/", "build/", ".pub-preload-cache/", ".ios/", ".android/",
		".macos/", "ephemeral/", "*.iml",
	}},
	"swift": {Key: "swift", Title: "Swift / Xcode", Category: CatLanguage, Order: 39, Rules: []string{
		".build/", "DerivedData/", "*.xcuserstate", "xcuserdata/", "*.moved-aside",
		"*.pbxuser", "*.perspectivev3", "*.hmap", "*.ipa", "*.dSYM.zip", "*.dSYM",
		"Pods/",
	}},
	"elixir": {Key: "elixir", Title: "Elixir / Mix", Category: CatLanguage, Order: 40, Rules: []string{
		"_build/", "deps/", ".elixir_ls/", "erl_crash.dump", ".fetch",
	}},
	"haskell": {Key: "haskell", Title: "Haskell / Stack", Category: CatLanguage, Order: 41, Rules: []string{
		"dist/", ".stack-work/", ".cabal-sandbox/", "*.hi", "*.o", ".hie/",
	}},
	"r-lang": {Key: "r-lang", Title: "R", Category: CatLanguage, Order: 42, Rules: []string{
		".Rhistory", ".Rapp.history", ".RData", ".Renviron", ".Ruserdata",
	}},

	// ---------------- Frameworks / build tools ----------------
	"next": {Key: "next", Title: "Next.js", Category: CatFramework, Order: 50, Rules: []string{
		".next/", "out/", ".pnp/", ".env*.local",
	}},
	"nuxt": {Key: "nuxt", Title: "Nuxt", Category: CatFramework, Order: 51, Rules: []string{
		".nuxt/", ".output/", ".data/", ".nitro/", ".nuxt/dist/",
	}},
	"svelte": {Key: "svelte", Title: "SvelteKit", Category: CatFramework, Order: 52, Rules: []string{
		".svelte-kit/", "build/", ".vercel",
	}},
	"astro": {Key: "astro", Title: "Astro", Category: CatFramework, Order: 53, Rules: []string{
		"dist/", ".astro/",
	}},
	"vite": {Key: "vite", Title: "Vite", Category: CatFramework, Order: 54, Rules: []string{
		"dist/", "node_modules/.vite/", "node_modules/.vite-temp/", "*.local",
	}},
	"angular": {Key: "angular", Title: "Angular", Category: CatFramework, Order: 55, Rules: []string{
		"dist/", ".angular/", "tmp/", ".sass-cache/",
	}},
	"gatsby": {Key: "gatsby", Title: "Gatsby", Category: CatFramework, Order: 56, Rules: []string{
		".cache/", "public/",
	}},
	"parcel": {Key: "parcel", Title: "Parcel", Category: CatFramework, Order: 57, Rules: []string{
		".parcel-cache/", "dist/",
	}},
	"turbo": {Key: "turbo", Title: "Turborepo", Category: CatFramework, Order: 58, Rules: []string{
		".turbo/",
	}},
	"tsc": {Key: "tsc", Title: "TypeScript / ESLint cache", Category: CatFramework, Order: 59, Rules: []string{
		"*.tsbuildinfo", ".eslintcache", ".stylelintcache",
	}},
	"storybook": {Key: "storybook", Title: "Storybook", Category: CatFramework, Order: 60, Rules: []string{
		"storybook-static/",
	}},
	"test-e2e": {Key: "test-e2e", Title: "Playwright / Cypress artifacts", Category: CatFramework, Order: 61, Rules: []string{
		"test-results/", "playwright-report/", "playwright/.cache/",
		"cypress/videos/", "cypress/screenshots/", "cypress/downloads/",
	}},
	"expo": {Key: "expo", Title: "Expo", Category: CatFramework, Order: 62, Rules: []string{
		".expo/", "dist/", "web-build/", "expo-env.d.ts",
	}},
	"jekyll": {Key: "jekyll", Title: "Jekyll / Hexo", Category: CatFramework, Order: 63, Rules: []string{
		"_site/", ".jekyll-cache/", ".jekyll-metadata", "public/", ".bundle/",
	}},
	"unity": {Key: "unity", Title: "Unity", Category: CatFramework, Order: 64, Rules: []string{
		"[Ll]ibrary/", "[Tt]emp/", "[Oo]bj/", "[Bb]uild/", "[Bb]uilds/",
		"[Ll]ogs/", "[Uu]ser[Ss]ettings/", "MemoryCaptures/", "[Uu]ser [Ii]nfo/",
	}},
	"unreal": {Key: "unreal", Title: "Unreal Engine", Category: CatFramework, Order: 65, Rules: []string{
		"Binaries/", "Intermediate/", "Saved/", "DerivedDataCache/", "Script/",
		"*.VC.db", "*.opensdf", "*.sdf",
	}},

	// ---------------- AI assistants / agents / MCP ----------------
	"claude-code": {Key: "claude-code", Title: "Claude Code", Category: CatAI, Order: 70, Rules: []string{
		".claude/settings.local.json",
		".claude/todos/",
		".claude/shell-snapshots/",
		".claude/statsig/",
		".claude/projects/",
		".claude/ide-logs/",
		".claude.json",
	}},
	"cursor": {Key: "cursor", Title: "Cursor", Category: CatAI, Order: 71, Rules: []string{
		".cursor/.tmp/",
		".cursor/eval/",
		".cursor/index/",
		".cursor/cache/",
		".cursor/*.local",
		".cursor/sessions/",
	}},
	"windsurf": {Key: "windsurf", Title: "Windsurf", Category: CatAI, Order: 72, Rules: []string{
		".windsurf/cache/",
		".windsurf/index/",
		".windsurf/.tmp/",
		".windsurf/*.local",
		".windsurf/sessions/",
	}},
	"codeium": {Key: "codeium", Title: "Codeium", Category: CatAI, Order: 73, Rules: []string{
		".codeium/", ".codeium.local/",
	}},
	"copilot": {Key: "copilot", Title: "GitHub Copilot", Category: CatAI, Order: 74, Rules: []string{
		".copilot/", ".github/copilot/local/",
	}},
	"aider": {Key: "aider", Title: "Aider", Category: CatAI, Order: 75, Rules: []string{
		".aider.input.history",
		".aider.llm.history",
		".aider.chat.history.md",
		".aider.conf.yml",
		".aider.tags.cache.v3/",
		".aider*",
	}},
	"continue": {Key: "continue", Title: "Continue", Category: CatAI, Order: 76, Rules: []string{
		".continue/",
	}},
	"cline": {Key: "cline", Title: "Cline", Category: CatAI, Order: 77, Rules: []string{
		".cline/", ".cline.local/",
	}},
	"roo": {Key: "roo", Title: "Roo Code", Category: CatAI, Order: 78, Rules: []string{
		".roo/", ".roo.local/",
	}},
	"cody": {Key: "cody", Title: "Sourcegraph Cody", Category: CatAI, Order: 79, Rules: []string{
		".cody/",
	}},
	"codex": {Key: "codex", Title: "OpenAI Codex CLI", Category: CatAI, Order: 80, Rules: []string{
		".codex/",
	}},
	"gemini-cli": {Key: "gemini-cli", Title: "Gemini CLI", Category: CatAI, Order: 81, Rules: []string{
		".gemini/",
	}},
	"opencode": {Key: "opencode", Title: "OpenCode", Category: CatAI, Order: 82, Rules: []string{
		".opencode/",
	}},
	"v0": {Key: "v0", Title: "v0", Category: CatAI, Order: 83, Rules: []string{
		".v0/",
	}},
	"mcp-local": {Key: "mcp-local", Title: "MCP local configs (may contain secrets)", Category: CatAI, Order: 84, Rules: []string{
		"mcp.local.json",
		".mcp.local.json",
		".mcp.json.local",
		".cursor/mcp.local.json",
		".vscode/mcp.local.json",
		"*.mcp.local",
	}},
	"tabnine": {Key: "tabnine", Title: "Tabnine", Category: CatAI, Order: 85, Rules: []string{
		".tabnine/",
	}},

	// ---------------- Environment & secrets ----------------
	"env": {Key: "env", Title: "Environment variables", Category: CatSecrets, Order: 90, Rules: []string{
		".env",
		".env.local",
		".env.*.local",
		".env.development.local",
		".env.test.local",
		".env.production.local",
		"!.env.example",
		"!.env.template",
		"!.env.sample",
	}},
	"cloud-keys": {Key: "cloud-keys", Title: "Cloud credentials & private keys", Category: CatSecrets, Order: 91, Rules: []string{
		".aws/credentials", ".gcp/", ".azure/credentials",
		"*.key", "*.p12", "*.pfx", "*.keystore", "*.jks",
		"id_rsa", "id_ed25519", "*.ppk",
		"service-account*.json", "*credentials*.json",
	}},
	"kubeconfig": {Key: "kubeconfig", Title: "Kubernetes kubeconfig", Category: CatSecrets, Order: 92, Rules: []string{
		"kubeconfig", "*.kubeconfig", ".kube/config",
	}},
	"terraform": {Key: "terraform", Title: "Terraform", Category: CatSecrets, Order: 93, Rules: []string{
		"*.tfstate", "*.tfstate.*", ".terraform/", "crash.log", "*.tfvars",
		"override.tf", "override.tf.json", "*.auto.tfvars",
	}},
	"ansible": {Key: "ansible", Title: "Ansible", Category: CatSecrets, Order: 94, Rules: []string{
		"*.retry", "host_vars/", "group_vars/",
	}},

	// ---------------- Deployment / cloud ----------------
	"docker": {Key: "docker", Title: "Docker", Category: CatDeploy, Order: 100, Rules: []string{
		"docker-compose.override.yml", ".docker/",
	}},
	"vercel": {Key: "vercel", Title: "Vercel", Category: CatDeploy, Order: 101, Rules: []string{
		".vercel",
	}},
	"netlify": {Key: "netlify", Title: "Netlify", Category: CatDeploy, Order: 102, Rules: []string{
		".netlify",
	}},
	"serverless": {Key: "serverless", Title: "Serverless / SST / Wrangler", Category: CatDeploy, Order: 103, Rules: []string{
		".serverless/", ".sst/", ".wrangler/",
	}},
	"firebase": {Key: "firebase", Title: "Firebase", Category: CatDeploy, Order: 104, Rules: []string{
		".firebase/", "firebase-debug.log", "ui-debug.log",
	}},

	// ---------------- Logs & temp ----------------
	"logs": {Key: "logs", Title: "Logs", Category: CatLogs, Order: 110, Rules: []string{
		"*.log", "logs/", "*.pid", "*.seed",
	}},
	"temp": {Key: "temp", Title: "Temporary & cache", Category: CatLogs, Order: 111, Rules: []string{
		".tmp/", "tmp/", "temp/", ".cache/",
	}},
}

// Get returns a template by key.
func Get(key string) (Template, bool) {
	t, ok := registry[key]
	return t, ok
}

// All returns every template.
func All() []Template {
	out := make([]Template, 0, len(registry))
	for _, t := range registry {
		out = append(out, t)
	}
	return out
}
