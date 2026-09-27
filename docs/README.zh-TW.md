<div align="center">

# 🦞 ignorekit

**離線、零相依的 `.gitignore` 產生器與 Git 追蹤健檢工具**

辨識技術棧 → 產生正確的 `.gitignore` → 找出你已經提交卻應該忽略的檔案 → 一道指令修好。

🌐 **切換語言**
[English](../README.md) · [简体中文](README.zh-CN.md) · [繁體中文](README.zh-TW.md)

⬇️ **[下載最新版本](https://github.com/gitstq/ignorekit/releases/latest)**

</div>

---

## 🎉 專案簡介

現在 AI 編碼助手與代理無所不在 —— **Claude Code、Cursor、Windsurf、GitHub
Copilot、Aider、Continue、Cline、Roo、Codex、MCP 伺服器**等等。它們會在你的儲存庫中
留下一連串檔案：工作階段紀錄、代理記憶、本機設定、快取，以及可能含有 **API 權杖與憑證**
的 MCP 設定。

一個不小心的 `git add .` 就可能把這些資料外洩。但同時，其中*有一部分*檔案原本就該和團隊
共用（例如 `.cursorrules`、`CLAUDE.md`、`copilot-instructions.md`）。要把兩者分辨開來既
繁瑣又容易出錯 —— 而傳統的 `.gitignore` 產生器對新興 AI 工具的涵蓋總是慢半拍，也從來不會
告訴你「**你早就提交了不該提交的東西**」。

**ignorekit** 正是為了解決這個問題。它是一個快速、單一執行檔、完全離線的工具，能夠：

- 🧠 理解現代語言、框架、編輯器以及 **AI 代理 / MCP**，清楚區分哪些檔案該**提交**、哪些
  該**忽略**；
- 🧪 把建議規則與 Git 已追蹤檔案交叉比對，揪出**已追蹤卻應忽略**的檔案，以及高把握度的
  密鑰警示；
- 🩹 把規則安全地整合進現有 `.gitignore`，並可用 `git rm --cached` 取消追蹤，過程具備冪等
  性，絕對不會更動你的工作樹檔案。

> 🍪 靈感來源：近期 GitHub Trending 上的代理工具浪潮（代理管理器、代理記憶、MCP 伺服器）
> 凸顯了這些工具會產生多少本機狀態。ignorekit 為獨立設計與實作，不複製任何現有專案；
> 忽略規則屬於通用的事實。

---

## ✨ 核心功能

- 🔌 **零執行階段相依** —— 單一靜態執行檔，以純 Go 撰寫（只使用標準函式庫）。不需要
  Node、Python、Docker，也不需要連網。
- 📴 **完全離線** —— 所有範本都已內建，任何資料都不會被上傳。
- 🤖 **AI 代理與 MCP 感知** —— 涵蓋 Claude Code、Cursor、Windsurf、Codeium、Copilot、
  Aider、Continue、Cline、Roo、Cody、Codex、Gemini CLI、OpenCode、Tabnine 以及本機
  MCP 設定，並明確指引「**保留/提交**」與「**忽略**」。
- 🧰 **技術棧涵蓋廣泛** —— 60+ 範本，橫跨作業系統、編輯器（VS Code、JetBrains、Vim、
  Emacs、Eclipse、Sublime）、語言（Node、Python、Go、Rust、Java、.NET、Ruby、PHP、
  Dart/Flutter、Swift、Elixir、Haskell、R）、框架（Next、Nuxt、SvelteKit、Astro、Vite、
  Angular、Gatsby、Parcel、Turborepo、Storybook、Playwright/Cypress、Expo、Unity、Unreal）
  等。
- 🧹 **追蹤健檢稽核** —— 找出 Git 已追蹤、卻符合忽略規則的檔案，這是傳統產生器欠缺的
  能力。
- 🔐 **密鑰警示** —— 一組少量、高把握度的憑證模式（GitHub/GitLab/AWS/Google/Slack/
  Stripe/OpenAI 權杖、JWT、私鑰）。它刻意不是完整的密鑰掃描器，請搭配專門工具使用。
- 🧷 **安全且冪等的整合** —— 保留你手寫的規則、自動移除重複、依章節分組，並在清楚標示的
  區塊中管理自身規則。
- 🛠️ **修復模式** —— 一步寫入檔案，並可選擇執行 `git rm --cached`；工作樹檔案一律保留。
- 🚦 **對 CI 友善** —— 提供 `text`、`json`、`SARIF 2.1.0` 輸出與穩定的結束代碼
  （0 乾淨，1 有發現，2 發生錯誤）。
- 🪟 **跨平台** —— Windows、macOS、Linux，支援 `amd64`/`arm64`（以及 `386`）。

---

## 🚀 快速開始

### 環境需求

除了執行檔本身，一無所需。已為 Windows、macOS、Linux 提供預先建置的執行檔。若要從原始碼
建置，則需要 **Go 1.23（含）以上**。

### 方式 A —— 下載預建置執行檔（建議）

1. 開啟 **[最新 Release](https://github.com/gitstq/ignorekit/releases/latest)**。
2. 下載符合你作業系統/處理器架構的壓縮檔。
3. （可選）驗證校驗碼，然後把 `ignorekit`（或 `ignorekit.exe`）放到 `PATH` 中。

```bash
# 以 Linux x86_64 為例
tar -xzf ignorekit_1.0.0_linux_amd64.tar.gz
sudo mv ignorekit_1.0.0_linux_amd64/ignorekit /usr/local/bin/
ignorekit version
```

### 方式 B —— 使用 Go 安裝

```bash
go install github.com/gitstq/ignorekit@latest
```

### 方式 C —— 從原始碼建置

```bash
git clone https://github.com/gitstq/ignorekit.git
cd ignorekit
make build          # 或：go build -trimpath -o ignorekit .
```

### 30 秒上手

```bash
cd your-project

# 1）檢視 ignorekit 將產生的 .gitignore（印到標準輸出）
ignorekit init

# 2）唯讀健檢：需要忽略的已追蹤檔案 + 密鑰警示
ignorekit audit

# 3）套用：把規則整合進 .gitignore，並取消追蹤應該忽略的檔案
ignorekit apply --untrack
```

---

## 📖 詳細使用指南

### 指令

| 指令 | 說明 |
|---|---|
| `init` | 偵測工具並產生 `.gitignore`（預設輸出到標準輸出） |
| `audit` | 唯讀健檢報告（已追蹤卻應忽略的檔案 + 密鑰） |
| `apply` | 把建議規則整合進 `.gitignore`（安全且冪等） |
| `list-templates` | 列出所有內建範本 |
| `version` | 印出版本資訊 |

### 全域旗標

| 旗標 | 說明 |
|---|---|
| `-C`、`--dir DIR` | 專案目錄（預設：目前目錄） |
| `-f`、`--format FMT` | 輸出格式：`text` \| `json` \| `sarif`（audit） |
| `--write` | 把結果寫入 `.gitignore`（init） |
| `--force` | 覆蓋已存在的 `.gitignore`（init） |
| `--untrack` | 對已追蹤卻應忽略的檔案一併執行 `git rm --cached`（apply） |
| `--dry-run` | 只顯示計畫，不寫入（apply） |
| `-y`、`--yes` | 不再詢問，直接進行 |
| `--no-color` | 關閉彩色輸出 |
| `--strict` | 當缺少 `.gitignore` 但有建議規則時，視為一項發現（audit） |

全域旗標放在指令**之前或之後**都可以，例如 `ignorekit -C ./app audit` 與
`ignorekit audit -C ./app` 皆有效。

### `init` —— 產生 `.gitignore`

```bash
ignorekit init                 # 在標準輸出預覽
ignorekit init --write         # 建立 .gitignore（若已存在會報錯）
ignorekit init --write --force # 覆蓋現有檔案
```

### `audit` —— 唯讀健檢報告

```bash
ignorekit audit                       # 人類可讀
ignorekit audit --format json         # 機器可讀
ignorekit audit --format sarif        # 供程式碼掃描使用的 SARIF 2.1.0
```

稽核會針對每個偵測到的工具，列出你應該**保留/提交**（綠色）與需要**檢視**（黃色）的檔案，
接著顯示：

- 符合忽略規則的已追蹤檔案（區分新規則 / 現有規則）；
- 高把握度密鑰警示，包含檔案與行號（已遮罩）。

### `apply` —— 整合並可選擇取消追蹤

```bash
ignorekit apply                 # 整合建議規則（建立受管理區塊）
ignorekit apply --dry-run       # 預覽將進行的變更
ignorekit apply --untrack       # 另外把已追蹤卻應忽略的檔案從索引移除
```

ignorekit 會把自身規則寫入一個邊界清楚的區塊：

```gitignore
# >>> ignorekit managed (safe to delete; regenerate with `ignorekit apply`) >>>
# --- Node.js ---
node_modules/
# <<< ignorekit managed <<<
```

區塊上方你手寫的規則絕不會被修改或重複；重複執行 `apply` 很安全：它會原地更新區塊，沒有
變更時顯示 `0 newly added`。

### AI 工具指引（該提交什麼 / 忽略什麼）

| 工具 | 提交 / 共用 | 忽略 |
|---|---|---|
| Claude Code | `CLAUDE.md`、`.claude/settings.json` | `.claude/settings.local.json`、todos、shell 快照、工作階段/專案狀態 |
| Cursor | `.cursorrules`、`.cursor/rules/` | 快取、索引、eval、工作階段、`*.local` |
| Windsurf | `.windsurf/rules/` | 快取、索引、工作階段、`*.local` |
| GitHub Copilot | `.github/copilot-instructions.md` | 本機 Copilot 狀態 |
| Aider | —— | 輸入/LLM/聊天紀錄、`.aider.conf.yml`、快取 |
| Continue / Cline / Roo / Cody / Codex | 規則檔（例如 `.clinerules`、`.roorules`、`AGENTS.md`） | 本機狀態目錄 |
| MCP | 經檢視後共用的伺服器設定 | `*.mcp.local`、`mcp.local.json` |

> 共用的 MCP 設定（`.mcp.json`、`mcp.json`）**不會被自動忽略**，因為團隊常會對它做版本
> 管控。ignorekit 改為掃描其中是否含有密鑰並提示檢視。只在本機使用的憑證請放入本機變體
> （`mcp.local.json`）。

### 持續整合（CI）

當存在追蹤健檢或密鑰發現時，`audit` 的結束代碼為 `1`；乾淨時為 `0`；出錯時為 `2` —— 非常
適合作為 CI 把關。

**GitHub Actions**（上傳 SARIF 至程式碼掃描）：

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

### 常見使用情境

- 🆕 **新專案** —— `ignorekit init --write` 取得完整基線。
- 🧑‍🤝‍🧑 **接手現有儲存庫** —— 先 `ignorekit audit`，再 `ignorekit apply --untrack`，
  檢視暫存的刪除後再提交。
- 🔒 **推送前的安全防線** —— 在推送前或 CI 中執行 `ignorekit audit`。
- 🤖 **安裝 AI 助手之後** —— 重新執行 `apply` 以涵蓋它的本機狀態。

### 截圖 / 展示

> 📸 預留位置 —— 可在此放入終端機錄影，例如 `![demo](demo.gif)`。

---

## 💡 設計理念與迭代規劃

### 設計原則

- **離線、本機優先。** 密鑰與原始碼永遠不應離開你的機器；範本內建，比對邏輯以純 Go 實作。
- **零相依。** 單一靜態執行檔可消除安裝門檻與供應鏈風險。
- **預設安全。** 所有會變更的指令都支援 dry-run，絕不刪除工作樹檔案，並把產生的規則隔離在
  可移除的區塊內。
- **高訊號、低雜訊。** 密鑰偵測以把握度優先，避免警示疲勞；全面掃描則交由專門工具。

### 為什麼選擇 Go？

Go 能以同一套程式碼產生小巧、靜態、跨平台的執行檔，擁有功能完善的標準函式庫（包含正規表達
式與旗標解析），並透過 `CGO_ENABLED=0` 讓所有發布目標的交叉編譯變得輕而易舉。

### 迭代規劃

- [ ] 涵蓋更多 AI 工具與框架簽章（由社群驅動）。
- [ ] 透過可選設定檔支援使用者自訂範本。
- [ ] 提供 `--interactive` 互動模式，逐項檢視發現。
- [ ] 安裝 pre-commit 鉤子，並提供更多 CI 範本。
- [ ] 在明確旗標之後提供更深入的密鑰/熵偵測。
- [ ] 補上日語、韓語、西班牙語 README。

非常歡迎透過 Issue 與 PR 貢獻功能、為路線圖投票 —— 詳見
[CONTRIBUTING.md](../CONTRIBUTING.md)。

---

## 📦 封裝與部署指南

ignorekit 屬於**桌面/命令列應用程式**，以預建置執行檔的形式發布。

### 建置全部發布壓縮檔

```bash
make cross          # 或：VERSION=1.0.0 ./scripts/build.sh
```

此指令會以 `CGO_ENABLED=0`、`-trimpath` 並移除符號，交叉編譯所有目標，然後把壓縮檔寫入
`dist/`：

- Windows：`.zip`（`amd64`、`arm64`、`386`）
- macOS / Linux：`.tar.gz`（`amd64`、`arm64`；Linux 另含 `386`）
- `dist/SHA256SUMS`：每個壓縮檔的 SHA-256 校驗碼

### 驗證下載

```bash
sha256sum -c SHA256SUMS --ignore-missing
```

### 相容環境

| 作業系統 | 處理器架構 |
|---|---|
| Windows | amd64、arm64、386 |
| macOS | amd64（Intel）、arm64（Apple Silicon） |
| Linux | amd64、arm64、386 |

所有執行檔皆為靜態連結，沒有任何執行階段程式庫相依。

---

## 🤝 貢獻指南

歡迎提交 PR 與 Issue！完整流程請閱讀 [CONTRIBUTING.md](../CONTRIBUTING.md)。簡單來說：

1. Fork 儲存庫，並從 `main` 切出主題分支。
2. 補上測試，確認 `go test ./...` 通過；執行 `gofmt -s` 與 `go vet`。
3. 使用[約定式提交](https://www.conventionalcommits.org/)
   （`feat:`、`fix:`、`docs:`、`refactor:`）。
4. 除非已在 Issue 中討論，否則請維持專案零相依。

請在所有互動中保持友善與建設性。

---

## 📄 授權條款說明

ignorekit 採用 **[MIT 授權條款](../LICENSE)** 釋出。

<div align="center">
🦞 為不想外洩代理密鑰的開發者而生。
</div>
