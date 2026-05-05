# DevMemory — 开发进度跟踪

> 每个步骤完成后更新此文档。状态：TODO / IN PROGRESS / DONE。

---

## 当前阶段：第一阶段完成

**开始日期**：2026-05-04
**目标**：最小可用个人记忆库 — 完成

---

## 步骤进度

### 1. 项目骨架

| # | 任务 | 状态 |
|---|------|------|
| 1.1 | 初始化 Go module | DONE |
| 1.2 | 创建目录结构 | DONE |

### 2. Core 模型

| # | 任务 | 状态 |
|---|------|------|
| 2.1 | `internal/core/types.go` — EntryType 常量 | DONE |
| 2.2 | `internal/core/entry.go` — Entry struct + 辅助方法 | DONE |

### 3. Store 层

| # | 任务 | 状态 |
|---|------|------|
| 3.1 | `internal/store/store.go` — Store interface + ListOptions | DONE |
| 3.2 | `internal/store/bbolt_store.go` — bbolt 实现 | DONE |
| 3.3 | `internal/store/bbolt_store_test.go` — 测试（11 cases） | DONE |

### 4. Config / Paths

| # | 任务 | 状态 |
|---|------|------|
| 4.1 | `internal/config/config.go` — Config struct | DONE |
| 4.2 | `internal/config/paths.go` — 通用逻辑 + portable mode | DONE |
| 4.3 | `internal/config/paths_linux.go` — Linux 路径 | DONE |
| 4.4 | `internal/config/paths_windows.go` — Windows 路径 | DONE |
| 4.5 | `internal/config/paths_darwin.go` — macOS 路径 | DONE |

### 5. Search

| # | 任务 | 状态 |
|---|------|------|
| 5.1 | `internal/search/search.go` — 搜索逻辑 | DONE |
| 5.2 | `internal/search/scorer.go` — 权重评分 | DONE |
| 5.3 | `internal/search/search_test.go` — 测试（7 cases） | DONE |

### 6. CLI（第 1 周）

| # | 任务 | 状态 |
|---|------|------|
| 6.1 | `cmd/devmemory/main.go` — 入口 + 子命令分发 | DONE |
| 6.2 | version 子命令 | DONE |
| 6.3 | add 子命令（含自动类型推断 + 危险检测） | DONE |
| 6.4 | list 子命令（type/project/tag 筛选） | DONE |
| 6.5 | search 子命令（多词 + 加权排序） | DONE |
| 6.6 | today 子命令（按日期 + 时间线格式） | DONE |

### 7. CLI（第 2 周）

| # | 任务 | 状态 |
|---|------|------|
| 7.1 | show 子命令（详情，支持短 ID 前缀匹配） | DONE |
| 7.2 | delete 子命令（确认提示 + --force 跳过） | DONE |
| 7.3 | edit 子命令（--title/--content/--type/--project/--tags/--favorite/--archive） | DONE |
| 7.4 | Store.ResolveID — 短 ID 前缀解析 | DONE |
| 7.5 | Windows 双击 exe 不闪退（waitForExit） | DONE |

### 8. Export / Import

| # | 任务 | 状态 |
|---|------|------|
| 8.1 | `internal/export/markdown.go` — Markdown daily export | DONE |
| 8.2 | `internal/export/json.go` — JSON export/import | DONE |
| 8.3 | `internal/export/export_test.go` — 测试（7 cases） | DONE |
| 8.4 | CLI export today（-o 输出到文件） | DONE |
| 8.5 | CLI export json（-o 输出到文件） | DONE |
| 8.6 | CLI import（ID 冲突时覆盖） | DONE |

### 9. 构建与文档

| # | 任务 | 状态 |
|---|------|------|
| 9.1 | `scripts/build.sh` — 四平台构建脚本 | DONE |
| 9.2 | `README.md` — 第一版 | DONE |

### 10. HTTP Server + Web UI（第 3 周）

| # | 任务 | 状态 |
|---|------|------|
| 10.1 | `internal/server/server.go` — HTTP 服务器 + embed 静态文件 | DONE |
| 10.2 | `internal/server/handlers.go` — REST API 全部 handler | DONE |
| 10.3 | `internal/server/exec.go` / `exec_windows.go` — 平台浏览器打开 | DONE |
| 10.4 | `internal/server/static/index.html` — Web UI 页面 | DONE |
| 10.5 | `internal/server/static/style.css` — 暗色主题样式 | DONE |
| 10.6 | `internal/server/static/app.js` — 前端交互逻辑 | DONE |
| 10.7 | CLI serve 子命令（--port / --no-open） | DONE |

---

## 集成验证

| # | 检查项 | 状态 |
|---|--------|------|
| V1 | `go build ./cmd/devmemory` 编译通过 | DONE |
| V2 | `go test ./...` 全部通过（25 tests） | DONE |
| V3 | CLI add → list 可工作 | DONE |
| V4 | CLI search 可返回结果 | DONE |
| V5 | CLI today 可列出当日记录 | DONE |
| V6 | 四平台交叉编译成功（linux/darwin/windows） | DONE |
| V7 | Windows 双击 exe 不再闪退 | DONE |
| V8 | CLI show 支持短 ID 前缀 | DONE |
| V9 | CLI delete 含确认提示，--force 跳过 | DONE |
| V10 | CLI edit 可更新各字段 | DONE |
| V11 | export today 输出 Markdown | DONE |
| V12 | export json / import JSON round-trip 数据无损 | DONE |
| V13 | HTTP server 启动并正确提供静态文件 | DONE |
| V14 | REST API 全部 endpoint 功能正确 | DONE |
| V15 | 短 ID 前缀匹配在 API 中工作 | DONE |
| V16 | 危险命令自动标记（dangerous flag） | DONE |
| V17 | Windows 交叉编译含 server 包成功 | DONE |
| V18 | Web UI Entry Edit 功能（模态框内编辑所有字段） | DONE |

---

## 已知问题

- 暂无

---

## 设计调整记录

| 日期 | 调整 |
|------|------|
| 2026-05-04 | CLI 使用手动参数解析（非 cobra），支持 flags 在任意位置，更符合用户习惯 |
| 2026-05-04 | go mod tidy 自动选择 bbolt v1.4.3，Go toolchain 更新为 1.24.11 |
| 2026-05-04 | Store 接口移除 Search 方法，搜索完全由 search.Engine 负责，职责更清晰 |
| 2026-05-04 | Store 新增 ResolveID 方法，支持短 ID 前缀匹配（无歧义时唯一） |
| 2026-05-04 | import 策略：ID 已存在时覆盖（先尝试 Create，失败则 Update） |

---

## 变更日志

| 日期 | 变更 |
|------|------|
| 2026-05-04 | 初始化项目文档（DESIGN.md, PLAN.md, SPEC.md, PROGRESS.md） |
| 2026-05-04 | 完成第 1 周全部开发：Core + Store + Search + Config + Action + CLI + 测试 + 构建 + README |
| 2026-05-04 | 修复 Windows 双击 exe 闪退问题 |
| 2026-05-04 | 完成第 2 周开发：show/delete/edit CLI + 短 ID 前缀匹配 + Markdown/JSON 导出导入 + 测试 |
| 2026-05-04 | 完成第 3 周开发：HTTP Server + REST API（13 endpoints）+ Web UI（暗色主题）+ serve CLI 命令 |
| 2026-05-05 | 完成第一阶段最后项：Web UI Entry Edit（模态框内编辑 title/content/type/project/tags） |
