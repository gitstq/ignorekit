<div align="center">

# 🦞 ignorekit

**离线、零依赖的 `.gitignore` 生成器与 Git 跟踪卫生工作台**

识别技术栈 → 生成正确的 `.gitignore` → 找出你已经提交却本该忽略的文件 → 一条命令修复。

🌐 **语言切换**
[English](../README.md) · [简体中文](README.zh-CN.md) · [繁體中文](README.zh-TW.md)

⬇️ **[下载最新版本](https://github.com/gitstq/ignorekit/releases/latest)**

</div>

---

## 🎉 项目介绍

如今，AI 编码助手与智能体无处不在 —— **Claude Code、Cursor、Windsurf、GitHub
Copilot、Aider、Continue、Cline、Roo、Codex、MCP 服务器**等等。它们会在你的仓库里
留下一连串文件：会话历史、智能体记忆、本地设置、缓存，以及可能包含 **API 密钥与凭证**
的 MCP 配置。

一次不小心的 `git add .` 就可能把这些数据泄露出去。但与此同时，其中*有一部分*文件
本来就该和团队共享（例如 `.cursorrules`、`CLAUDE.md`、`copilot-instructions.md`）。
把两者区分开既繁琐又容易出错 —— 而传统的 `.gitignore` 生成器对新兴 AI 工具覆盖滞后，
更不会告诉你「**你已经提交了本不该提交的东西**」。

**ignorekit** 正是为此而生。它是一个快速、单二进制、完全离线的工具，能够：

- 🧠 理解现代语言、框架、编辑器以及 **AI 智能体 / MCP**，明确区分哪些文件该**提交**、
  哪些该**忽略**；
- 🧪 将建议规则与 Git 已跟踪文件交叉比对，暴露**已跟踪却应忽略**的文件，以及高置信度的
  密钥红旗；
- 🩹 把规则安全地合并进现有 `.gitignore`，并可用 `git rm --cached` 取消跟踪，过程幂等、
  绝不改动你的工作区文件。

> 🍪 灵感来源：近期 GitHub Trending 上的智能体工具浪潮（智能体管理器、智能体记忆、
> MCP 服务器）让我们看到这些工具会产生多少本地状态。ignorekit 为独立设计与实现，
> 不复制任何现有项目；忽略规则属于通用事实。

---

## ✨ 核心特性

- 🔌 **零运行时依赖** —— 单个静态二进制，纯 Go 编写（仅用标准库）。无需 Node、Python、
  Docker，也无需联网。
- 📴 **完全离线** —— 所有模板均已内嵌，任何数据都不会被上传。
- 🤖 **AI 智能体与 MCP 感知** —— 覆盖 Claude Code、Cursor、Windsurf、Codeium、Copilot、
  Aider、Continue、Cline、Roo、Cody、Codex、Gemini CLI、OpenCode、Tabnine 以及本地
  MCP 配置，并给出明确的「**保留/提交** vs **忽略**」指引。
- 🧰 **技术栈覆盖广** —— 60+ 模板，横跨操作系统、编辑器（VS Code、JetBrains、Vim、
  Emacs、Eclipse、Sublime）、语言（Node、Python、Go、Rust、Java、.NET、Ruby、PHP、
  Dart/Flutter、Swift、Elixir、Haskell、R）、框架（Next、Nuxt、SvelteKit、Astro、Vite、
  Angular、Gatsby、Parcel、Turborepo、Storybook、Playwright/Cypress、Expo、Unity、Unreal）
  等。
- 🧹 **跟踪卫生审计** —— 找出 Git 已跟踪、却命中忽略规则的文件，这是传统生成器不具备的
  能力。
- 🔐 **密钥红旗** —— 一组少量、高置信度的凭证模式（GitHub/GitLab/AWS/Google/Slack/
  Stripe/OpenAI 令牌、JWT、私钥）。它刻意不是完整的密钥扫描器，请配合专用工具使用。
- 🧷 **安全幂等合并** —— 保留你的手写规则、自动去重、按节分组，并在清晰标记的区块内管理
  自有规则。
- 🛠️ **修复模式** —— 一步写入文件，并可选执行 `git rm --cached`；工作区文件始终保留。
- 🚦 **对 CI 友好** —— 支持 `text`、`json`、`SARIF 2.1.0` 输出与稳定退出码
  （0 干净，1 有发现，2 出错）。
- 🪟 **跨平台** —— Windows、macOS、Linux，支持 `amd64`/`arm64`（以及 `386`）。

---

## 🚀 快速开始

### 环境要求

除了二进制本身，无需任何依赖。已为 Windows、macOS、Linux 提供预编译可执行文件。若从源码
构建，则需要 **Go 1.23+**。

### 方式 A —— 下载预编译二进制（推荐）

1. 打开 **[最新 Release](https://github.com/gitstq/ignorekit/releases/latest)**。
2. 下载对应操作系统/架构的压缩包。
3. （可选）校验完整性，并把 `ignorekit`（或 `ignorekit.exe`）放到 `PATH` 中。

```bash
# 以 Linux x86_64 为例
tar -xzf ignorekit_1.0.0_linux_amd64.tar.gz
sudo mv ignorekit_1.0.0_linux_amd64/ignorekit /usr/local/bin/
ignorekit version
```

### 方式 B —— 使用 Go 安装

```bash
go install github.com/gitstq/ignorekit@latest
```

### 方式 C —— 从源码构建

```bash
git clone https://github.com/gitstq/ignorekit.git
cd ignorekit
make build          # 或：go build -trimpath -o ignorekit .
```

### 30 秒上手

```bash
cd your-project

# 1）查看 ignorekit 将生成的 .gitignore（打印到标准输出）
ignorekit init

# 2）只读卫生审计：需要忽略的已跟踪文件 + 密钥红旗
ignorekit audit

# 3）应用：把规则合并进 .gitignore，并取消跟踪本该忽略的文件
ignorekit apply --untrack
```

---

## 📖 详细使用指南

### 子命令

| 命令 | 说明 |
|---|---|
| `init` | 检测工具并生成 `.gitignore`（默认输出到标准输出） |
| `audit` | 只读卫生报告（已跟踪却应忽略的文件 + 密钥） |
| `apply` | 把建议规则合并进 `.gitignore`（安全且幂等） |
| `list-templates` | 列出所有内置模板 |
| `version` | 打印版本信息 |

### 全局参数

| 参数 | 说明 |
|---|---|
| `-C`、`--dir DIR` | 项目目录（默认：当前目录） |
| `-f`、`--format FMT` | 输出格式：`text` \| `json` \| `sarif`（audit） |
| `--write` | 将结果写入 `.gitignore`（init） |
| `--force` | 覆盖已存在的 `.gitignore`（init） |
| `--untrack` | 对已跟踪却应忽略的文件一并执行 `git rm --cached`（apply） |
| `--dry-run` | 只展示计划，不写入（apply） |
| `-y`、`--yes` | 不再询问，直接执行 |
| `--no-color` | 关闭彩色输出 |
| `--strict` | 当缺少 `.gitignore` 但存在建议规则时，视为发现（audit） |

全局参数放在命令**之前或之后**均可，例如 `ignorekit -C ./app audit` 与
`ignorekit audit -C ./app` 都有效。

### `init` —— 生成 `.gitignore`

```bash
ignorekit init                 # 在标准输出预览
ignorekit init --write         # 创建 .gitignore（若已存在则报错）
ignorekit init --write --force # 覆盖已有文件
```

### `audit` —— 只读卫生报告

```bash
ignorekit audit                       # 人类可读
ignorekit audit --format json         # 机器可读
ignorekit audit --format sarif        # 用于代码扫描的 SARIF 2.1.0
```

审计会针对每个检测到的工具，列出你应当**保留/提交**（绿色）与需要**复核**（黄色）的文件，
随后给出：

- 命中忽略规则的已跟踪文件（区分新规则 / 现有规则）；
- 高置信度密钥红旗，含文件与行号（已脱敏）。

### `apply` —— 合并并可选取消跟踪

```bash
ignorekit apply                 # 合并建议规则（创建受管区块）
ignorekit apply --dry-run       # 预览将要发生的变更
ignorekit apply --untrack       # 额外把已跟踪却应忽略的文件从索引中移除
```

ignorekit 会把自有规则写入一个边界清晰的区块：

```gitignore
# >>> ignorekit managed (safe to delete; regenerate with `ignorekit apply`) >>>
# --- Node.js ---
node_modules/
# <<< ignorekit managed <<<
```

区块上方你手写的规则绝不会被修改或重复；重复运行 `apply` 是安全的：它会原地更新区块，
没有变更时报告 `0 newly added`。

### AI 工具指引（提交什么 / 忽略什么）

| 工具 | 提交 / 共享 | 忽略 |
|---|---|---|
| Claude Code | `CLAUDE.md`、`.claude/settings.json` | `.claude/settings.local.json`、todos、shell 快照、会话/项目状态 |
| Cursor | `.cursorrules`、`.cursor/rules/` | 缓存、索引、eval、会话、`*.local` |
| Windsurf | `.windsurf/rules/` | 缓存、索引、会话、`*.local` |
| GitHub Copilot | `.github/copilot-instructions.md` | 本地 Copilot 状态 |
| Aider | —— | 输入/LLM/聊天历史、`.aider.conf.yml`、缓存 |
| Continue / Cline / Roo / Cody / Codex | 规则文件（如 `.clinerules`、`.roorules`、`AGENTS.md`） | 本地状态目录 |
| MCP | 经复核后共享的服务器配置 | `*.mcp.local`、`mcp.local.json` |

> 共享的 MCP 配置（`.mcp.json`、`mcp.json`）**不会被自动忽略**，因为团队常会对其做版本
> 管理。ignorekit 改为扫描其中是否含密钥并提示复核。仅本机使用的凭证请放入本地变体
> （`mcp.local.json`）。

### 持续集成（CI）

当存在跟踪卫生或密钥发现时，`audit` 退出码为 `1`；干净时为 `0`；出错时为 `2` —— 非常
适合作为 CI 门禁。

**GitHub Actions**（上传 SARIF 至代码扫描）：

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

### 典型使用场景

- 🆕 **新项目** —— `ignorekit init --write` 生成完整基线。
- 🧑‍🤝‍🧑 **接手已有仓库** —— 先 `ignorekit audit`，再 `ignorekit apply --untrack`，
  检查暂存的删除操作后提交。
- 🔒 **推送前的安全网** —— 在推送前或 CI 中运行 `ignorekit audit`。
- 🤖 **安装 AI 助手之后** —— 重新运行 `apply` 以覆盖其本地状态。

### 截图 / 演示

> 📸 占位 —— 可在此放置终端录屏，例如 `![demo](demo.gif)`。

---

## 💡 设计思路与迭代规划

### 设计理念

- **离线、本地优先。** 密钥与源码永远不应离开你的机器；模板内嵌，匹配逻辑用纯 Go 实现。
- **零依赖。** 单个静态二进制消除安装门槛与供应链风险。
- **默认安全。** 所有变更类命令都支持 dry-run，绝不删除工作区文件，并把生成规则隔离在可
  删除的区块内。
- **高信号、低噪声。** 密钥检测优先保证置信度，避免告警疲劳；全面扫描交给专用工具。

### 为什么选择 Go？

Go 能在同一套代码下产出小巧、静态、跨平台的二进制，拥有功能完善的标准库（含正则与参数
解析），并通过 `CGO_ENABLED=0` 让所有发布目标的交叉编译变得轻而易举。

### 迭代规划

- [ ] 覆盖更多 AI 工具与框架签名（由社区驱动）。
- [ ] 通过可选配置文件支持用户自定义模板。
- [ ] 提供 `--interactive` 交互模式，逐条复核发现。
- [ ] 安装 pre-commit 钩子，并提供更多 CI 模板。
- [ ] 在显式参数后提供更深入的密钥/熵检测。
- [ ] 补充日语、韩语、西班牙语 README。

非常欢迎通过 Issue 与 PR 贡献功能、为路线图投票 —— 详见
[CONTRIBUTING.md](../CONTRIBUTING.md)。

---

## 📦 打包与部署指南

ignorekit 属于**桌面/命令行应用**，以预编译二进制形式发布。

### 构建全部发布压缩包

```bash
make cross          # 或：VERSION=1.0.0 ./scripts/build.sh
```

该命令以 `CGO_ENABLED=0`、`-trimpath` 并去除符号表交叉编译所有目标，随后把压缩包写入
`dist/`：

- Windows：`.zip`（`amd64`、`arm64`、`386`）
- macOS / Linux：`.tar.gz`（`amd64`、`arm64`；Linux 另含 `386`）
- `dist/SHA256SUMS`：每个压缩包的 SHA-256 校验和

### 校验下载

```bash
sha256sum -c SHA256SUMS --ignore-missing
```

### 兼容环境

| 操作系统 | 架构 |
|---|---|
| Windows | amd64、arm64、386 |
| macOS | amd64（Intel）、arm64（Apple Silicon） |
| Linux | amd64、arm64、386 |

所有二进制均为静态链接，无任何运行时库依赖。

---

## 🤝 贡献指南

欢迎提交 PR 与 Issue！完整流程请阅读 [CONTRIBUTING.md](../CONTRIBUTING.md)。简而言之：

1. Fork 仓库，并从 `main` 切出分支。
2. 补充测试，确保 `go test ./...` 通过；运行 `gofmt -s` 与 `go vet`。
3. 遵循[约定式提交](https://www.conventionalcommits.org/)
   （`feat:`、`fix:`、`docs:`、`refactor:`）。
4. 除非已在 Issue 中讨论，否则请保持项目零依赖。

请在所有交流中保持友善与建设性。

---

## 📄 开源协议说明

ignorekit 基于 **[MIT 协议](../LICENSE)** 发布。

<div align="center">
🦞 为不愿泄露智能体密钥的开发者而生。
</div>
