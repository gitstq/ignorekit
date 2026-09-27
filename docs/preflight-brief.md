# ignorekit 项目预研简报（2026-09-27）

## 1. 榜单观察（灵感来源）
当日 GitHub Trending（日/周榜）被 AI Coding Agent 生态主导：
`paperclipai/paperclip`（管理工作中的 Agent）、`vectorize-io/hindsight`（Agent 记忆）、
`anthropics/claude-code`、`addyosmani/agent-skills`、`mobile-next/mobile-mcp`、
`alibaba/open-code-review`、`Tencent/WeKnora` 等。

这些 Agent / MCP / 助手工具会在**项目仓库内**产生大量本地产物：会话历史、记忆、
本地设置、缓存、以及可能内嵌密钥的 MCP 配置。真实痛点：
- 这些文件极易被 `git add .` 误提交（会话日志、本地配置、密钥泄露）；
- 但其中又有一部分**应当共享**（如 `.cursorrules`、`.claude/settings.json`、
  `copilot-instructions.md`），「该忽略哪些 / 该提交哪些」缺乏明确工具；
- 传统 gitignore 模板（gitignore.io 等）离线能力弱、对 AI 工具覆盖滞后，
  且不会告诉你「你已经把不该提交的文件提交了」。

## 2. 自研定位（差异化，非复制）
**ignorekit**：零运行时依赖、单二进制、离线优先的现代 `.gitignore` 生成器
与 **Git 跟踪卫生（tracking-hygiene）工作台**。

不复制任何现有项目代码；gitignore 模式为通用事实规则，全部模板与引擎独立编写。

## 3. 核心差异化亮点
1. **AI Agent / MCP 感知**：内置 Claude Code、Cursor、Windsurf、Codeium、Copilot、
   Aider、Continue、Cline/Roo、Cody、Codex、Tabnine、MCP 等现代工具签名，
   并区分「应共享」与「应忽略」文件。
2. **跟踪卫生审计**：与 `git ls-files` 交叉比对，找出**已被跟踪但应忽略**的文件
   （含 Agent 状态/密钥），这是传统模板生成器不具备的能力。
3. **安全合并**：解析现有 `.gitignore`，去重、保留手写规则、分节合并、识别冲突否定。
4. **修复模式**：一键写入并可选 `git rm --cached` 取消跟踪（默认 dry-run）。
5. **高信号密钥红旗**：对已跟踪文件做少量高置信度密钥前缀检测（明确不是完整密钥扫描器）。
6. **CI 友好**：text / json / sarif 输出与退出码门禁；全平台单二进制。

## 4. 技术栈与架构
- 语言：Go 1.23（仅标准库，零第三方依赖），单二进制，便于全平台交叉编译。
- 模块：detect（检测）/ templates（内嵌模板, go:embed）/ gitignore（解析·匹配·合并）/
  gitx（git 集成）/ secret（密钥红旗）/ report（多格式输出）/ term（彩色文本）。
- 子命令：`init`（生成/写入）、`audit`（只读卫生审计）、`apply`（合并+可选取消跟踪）、
  `list-templates`、`version`。

## 5. 自测标准
- 单元测试：gitignore 解析与匹配（*、**、!、目录、锚定、字符类）、合并去重、检测签名、报告。
- 集成测试：临时 git 仓库中提交应忽略文件 → audit 报出 → apply --untrack → 验证取消跟踪。
- 全平台交叉编译产物 + SHA256 校验和；本地运行 audit/apply 全流程通过。

## 6. 相似度校验结论
账号现有 1688 个仓库中：gitignore 生成 0、跟踪卫生/untrack 0、ignore 0；
最接近的为 SecretShield-CLI（密钥扫描）与 BranchSweeper（分支清理），
功能定位、核心场景相似度均显著低于 60%，判定为非重复项目，可进入开发。
