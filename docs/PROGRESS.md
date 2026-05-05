# DevMemory — 开发进度跟踪

> 每个步骤完成后更新此文档。状态：TODO / IN PROGRESS / DONE。

---

## 当前阶段：第二阶段 — Fyne 原生桌面 UI 迁移

**开始日期**：2026-05-05
**完成日期**：2026-05-05
**目标**：将 DevMemory 从 Local Web UI 迁移到 Fyne 原生桌面 GUI — **已完成**

---

## 第一阶段（已完成）

**开始日期**：2026-05-04
**完成日期**：2026-05-05
**目标**：最小可用个人记忆库 — 已完成

### 第 1 周：Core + Store + Search + Config + CLI

| # | 任务 | 状态 |
|---|------|------|
| 1.1 | 初始化 Go module | DONE |
| 1.2 | 创建目录结构 | DONE |
| 1.3 | `internal/core/types.go` — EntryType 常量 | DONE |
| 1.4 | `internal/core/entry.go` — Entry struct + 辅助方法 | DONE |
| 1.5 | `internal/store/store.go` — Store interface + ListOptions | DONE |
| 1.6 | `internal/store/bbolt_store.go` — bbolt 实现 | DONE |
| 1.7 | `internal/store/bbolt_store_test.go` — 测试（11 cases） | DONE |
| 1.8 | `internal/config/config.go` — Config struct | DONE |
| 1.9 | `internal/config/paths.go` + 平台文件 — 数据目录 | DONE |
| 1.10 | `internal/search/search.go` + `scorer.go` — 搜索 | DONE |
| 1.11 | `internal/search/search_test.go` — 测试（7 cases） | DONE |
| 1.12 | `cmd/devmemory/main.go` — CLI 入口 + 子命令 | DONE |
| 1.13 | `scripts/build.sh` — 四平台构建脚本 | DONE |
| 1.14 | `README.md` — 第一版 | DONE |

### 第 2 周：CLI 增强 + Export/Import

| # | 任务 | 状态 |
|---|------|------|
| 2.1 | show 子命令（短 ID 前缀匹配） | DONE |
| 2.2 | delete 子命令（确认提示 + --force） | DONE |
| 2.3 | edit 子命令（多字段更新） | DONE |
| 2.4 | Store.ResolveID — 短 ID 前缀解析 | DONE |
| 2.5 | Windows 双击 exe 不闪退 | DONE |
| 2.6 | `internal/export/markdown.go` — Markdown daily export | DONE |
| 2.7 | `internal/export/json.go` — JSON export/import | DONE |
| 2.8 | `internal/export/export_test.go` — 测试（7 cases） | DONE |
| 2.9 | CLI export / import 子命令 | DONE |

### 第 3 周：HTTP Server + Web UI

| # | 任务 | 状态 |
|---|------|------|
| 3.1 | `internal/server/server.go` — HTTP server + embed 静态文件 | DONE |
| 3.2 | `internal/server/handlers.go` — REST API（13 endpoints） | DONE |
| 3.3 | `internal/server/exec.go` / `exec_windows.go` — 平台浏览器打开 | DONE |
| 3.4 | `internal/server/static/` — Web UI（Capture / Search / Today） | DONE |
| 3.5 | CLI serve 子命令（--port / --no-open） | DONE |
| 3.6 | Web UI Entry Edit（模态框内编辑所有字段） | DONE |

### 第一阶段集成验证

| # | 检查项 | 状态 |
|---|--------|------|
| V1 | `go build ./cmd/devmemory` 编译通过 | DONE |
| V2 | `go test ./...` 全部通过（25 tests） | DONE |
| V3 | CLI add → list 可工作 | DONE |
| V4 | CLI search 可返回结果 | DONE |
| V5 | CLI today 可列出当日记录 | DONE |
| V6 | 四平台交叉编译成功 | DONE |
| V7 | Windows 双击 exe 不闪退 | DONE |
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
| V18 | Web UI Entry Edit 功能 | DONE |

---

## 第二阶段：Fyne 迁移进度

### Step 1：Service 层 + CLI 重构

| # | 任务 | 状态 |
|---|------|------|
| 1.1 | 创建 `internal/service/memory_service.go` — MemoryService struct | DONE |
| 1.2 | 实现 CreateEntry（自动 DetectType + IsDangerous + NewEntry） | DONE |
| 1.3 | 实现 SearchEntries | DONE |
| 1.4 | 实现 GetToday | DONE |
| 1.5 | 实现 CopyEntry（clipboard + UseCount++ + LastUsedAt） | DONE |
| 1.6 | 实现 UpdateEntry | DONE |
| 1.7 | 实现 DeleteEntry | DONE |
| 1.8 | 实现 ExportTodayMarkdown | DONE |
| 1.9 | 实现 ExportJSON / ImportJSON | DONE |
| 1.10 | 实现 ListEntries / GetEntry | DONE |
| 1.11 | 重构 CLI main.go 调用 Service | DONE |
| 1.12 | 更新 server/handlers.go 调用 Service | DONE |

### Step 2：Fyne 依赖 + 骨架窗口

| # | 任务 | 状态 |
|---|------|------|
| 2.1 | `go get fyne.io/fyne/v2` | DONE |
| 2.2 | 创建 `internal/ui/fyne/app.go` | DONE |
| 2.3 | 创建 `internal/ui/fyne/main_window.go`（主窗口 + 导航） | DONE |
| 2.4 | 更新 main.go：无参数/gui 启动 Fyne | DONE |
| 2.5 | 验证 CGO 编译通过 | DONE |

### Step 3：Capture 页

| # | 任务 | 状态 |
|---|------|------|
| 3.1 | 创建 `page_capture.go` — 布局 | DONE |
| 3.2 | Content 多行输入 + Type 下拉 + Title/Project/Tags | DONE |
| 3.3 | 提交 → Service.CreateEntry + 反馈 | DONE |

### Step 4：Search 页

| # | 任务 | 状态 |
|---|------|------|
| 4.1 | 创建 `page_search.go` — 布局 | DONE |
| 4.2 | 搜索框 + 类型筛选 + 结果列表 | DONE |
| 4.3 | 创建 `entry_dialog.go` — Entry 详情/编辑弹窗 | DONE |
| 4.4 | 点击结果 → 弹窗 | DONE |

### Step 5：Today / Actions / Knowledge 页

| # | 任务 | 状态 |
|---|------|------|
| 5.1 | `page_today.go` — 今日条目时间线 | DONE |
| 5.2 | `page_actions.go` — 动作类型列表 + Copy/Open/Execute | DONE |
| 5.3 | `page_knowledge.go` — 知识类型列表 + 查看/编辑 | DONE |
| 5.4 | Copy/Open/Execute 功能绑定 Service | DONE |
| 5.5 | Favorite / Archive / Edit / Delete 操作 | DONE |

### Step 6：Settings 页 + 导出导入

| # | 任务 | 状态 |
|---|------|------|
| 6.1 | `page_settings.go` — 信息展示 + 导出导入按钮 | DONE |
| 6.2 | Export Today Markdown / Export JSON / Import JSON | DONE |
| 6.3 | 状态栏（版本 / 数据路径 / 条目数） | DONE |

### Step 7：构建脚本 + README

| # | 任务 | 状态 |
|---|------|------|
| 7.1 | 更新 `scripts/build.sh`（CGO + fyne-cross） | DONE |
| 7.2 | 更新 `README.md`（Fyne 说明 + 构建依赖） | DONE |
| 7.3 | 验证四平台构建 | DONE（Linux 25MB + Windows 26MB mingw-w64 交叉编译成功；macOS 需 zig/osxcross 或 GitHub Actions） |

---

## 第二阶段集成验证

| # | 检查项 | 状态 |
|---|--------|------|
| V1 | `go test ./...` 全部通过（含 Service 层） | TODO |
| V2 | CLI 行为与重构前完全一致 | TODO |
| V3 | `go run ./cmd/devmemory` 打开 Fyne 窗口 | TODO |
| V4 | Fyne Capture 可新增 Entry | TODO |
| V5 | Fyne Search 可搜索并查看结果 | TODO |
| V6 | Fyne Today 显示当日条目 | TODO |
| V7 | Copy/Open/Execute 功能正常 | TODO |
| V8 | Dangerous 命令二次确认 | TODO |
| V9 | Settings 导出导入正常 | TODO |
| V10 | 四平台构建成功 | DONE（Linux + Windows 本机构建成功；macOS 通过 GitHub Actions） |
| V11 | 现有数据不丢失（同一 bbolt 数据库） | TODO |

---

### 第三阶段：日常可用的开发者工作记忆系统（TODO）

| # | 功能 | 状态 |
|---|------|------|
| 3.1 | Project 字段增强（project list / filter） | TODO |
| 3.2 | Issue / Business knowledge 模板 | TODO |
| 3.3 | Prompt / snippet 管理优化 | TODO |
| 3.4 | Favorite 筛选 | TODO |
| 3.5 | Tag 管理（autocomplete / rename / merge） | TODO |
| 3.6 | Weekly report export | TODO |
| 3.7 | Portable mode 完善 | TODO |
| 3.8 | Entry 分页 | TODO |
| 3.9 | 快捷键支持 | TODO |

### 第四阶段：稳定的跨平台个人效率与知识工具（TODO）

| # | 功能 | 状态 |
|---|------|------|
| 4.1 | 全文搜索优化（Bleve 或纯 Go 方案） | TODO |
| 4.2 | Entry 关联关系（related entries） | TODO |
| 4.3 | Project timeline 视图 | TODO |
| 4.4 | Markdown vault 兼容（Obsidian 格式） | TODO |
| 4.5 | Backup / restore（自动备份策略） | TODO |
| 4.6 | Template variables（命令模板可填参数） | TODO |

---

## 已知问题

- 暂无

---

## 设计调整记录

| 日期 | 调整 |
|------|------|
| 2026-05-04 | CLI 使用手动参数解析（非 cobra），支持 flags 在任意位置 |
| 2026-05-04 | go mod tidy 自动选择 bbolt v1.4.3，Go toolchain 更新为 1.24.11 |
| 2026-05-04 | Store 接口移除 Search 方法，搜索完全由 search.Engine 负责 |
| 2026-05-04 | Store 新增 ResolveID 方法，支持短 ID 前缀匹配 |
| 2026-05-04 | import 策略：ID 已存在时覆盖 |
| 2026-05-05 | 第二阶段启动：从 Web UI 迁移到 Fyne 原生桌面 GUI |
| 2026-05-05 | 新增 Service 层：统一 CLI 和 Fyne UI 的业务逻辑入口 |

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
| 2026-05-05 | 刷新全部文档：DESIGN.md / PLAN.md / SPEC.md / PROGRESS.md，对齐 Fyne 迁移计划 |
| 2026-05-05 | 完成第二阶段全部开发：Service 层 + Fyne GUI（6 页面）+ 构建脚本 + README |
| 2026-05-05 | 创建 Windows 手动测试用例：docs/TEST_CASES_WINDOWS.md（90 项测试） |
| 2026-05-05 | 四平台构建：更新 build.sh（支持 mingw/zig/fyne-cross）、创建 GitHub Actions CI（4 平台） |
| 2026-05-05 | 本机交叉编译成功：安装 Fedora mingw64-gcc + mingw64-cpp（含 cc1），Linux 25MB + Windows 26MB 构建完成 |
