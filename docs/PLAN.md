# DevMemory — 路线图与任务拆解

---

## 阶段概览

| 阶段 | 目标 | 状态 |
|------|------|------|
| **第一阶段** | 最小可用个人记忆库（Web UI + CLI） | **DONE** |
| **第二阶段** | Fyne 原生桌面 UI 迁移（插入） | TODO |
| **第三阶段** | 日常可用的开发者工作记忆系统 | TODO |
| **第四阶段** | 稳定的跨平台个人效率与知识工具 | TODO |

---

## 第一阶段回顾（已完成）

**开始日期**：2026-05-04
**完成日期**：2026-05-05

**成果**：
- CLI 全部子命令（add/list/search/today/show/delete/edit/export/import/version）
- HTTP Server + REST API（13 endpoints）
- Web UI（Capture / Search / Today + Entry 详情编辑模态框）
- bbolt 存储 + 加权搜索 + Markdown/JSON 导出导入
- 四平台交叉编译

---

## 第二阶段：Fyne 原生桌面 UI 迁移

### 迁移策略

增量迁移，不重写。核心原则：

1. **先抽取 Service 层**：CLI 和 Fyne UI 共用同一套业务逻辑
2. **CLI 不变**：重构为调用 Service，行为完全一致
3. **Web UI 保留但降级**：`internal/server/` 代码不动，降级为备用入口
4. **Fyne UI 只调用 Service**：widget 回调中不含业务逻辑
5. **数据不丢失**：同一 bbolt 数据库，同一数据目录

### 7 步迁移计划

---

#### Step 1：抽取 Service 层 + CLI 重构

**目标**：CLI 和未来 Fyne UI 共用同一套业务逻辑，消除 main.go 和 handlers.go 中的重复代码。

**任务**：

| # | 任务 | 说明 |
|---|------|------|
| 1.1 | 创建 `internal/service/memory_service.go` | MemoryService struct，持有 Store + SearchEngine + Executor |
| 1.2 | 实现 `CreateEntry(input)` | 自动 DetectType + IsDangerous + NewEntry + store.Create |
| 1.3 | 实现 `SearchEntries(query, type)` | 调用 search.Engine |
| 1.4 | 实现 `GetToday()` | 调用 store.GetByDate |
| 1.5 | 实现 `CopyEntry(id)` | clipboard + UseCount++ + LastUsedAt |
| 1.6 | 实现 `UpdateEntry(id, fields)` | store.Get + 修改字段 + store.Update |
| 1.7 | 实现 `DeleteEntry(id)` | store.Delete |
| 1.8 | 实现 `ExportTodayMarkdown()` | 调用 export.ExportDailyMarkdown |
| 1.9 | 实现 `ExportJSON()` / `ImportJSON()` | 调用 export 包 |
| 1.10 | 实现 `ListEntries(opts)` | store.List |
| 1.11 | 实现 `GetEntry(id)` | ResolveID + store.Get |
| 1.12 | 重构 `cmd/devmemory/main.go` | 所有 CLI 子命令改为调用 Service，删除内联业务逻辑 |
| 1.13 | 更新 `internal/server/handlers.go` | handler 调用 Service 而非直接操作 Store（可选，降级优先级） |

**验收**：
- `go test ./...` 全部通过
- CLI 行为与重构前完全一致
- main.go 中不再有 DetectType / IsDangerous 等业务逻辑

---

#### Step 2：添加 Fyne 依赖 + 骨架窗口

**目标**：`devmemory` 或 `devmemory gui` 能打开一个 Fyne 窗口。

**任务**：

| # | 任务 | 说明 |
|---|------|------|
| 2.1 | `go get fyne.io/fyne/v2` | 添加 Fyne 依赖 |
| 2.2 | 创建 `internal/ui/fyne/app.go` | Fyne app 初始化 + Service 绑定 |
| 2.3 | 创建 `internal/ui/fyne/main_window.go` | 主窗口 + 左侧导航（6 个 tab） |
| 2.4 | 更新 `cmd/devmemory/main.go` | 无参数或 `gui` 子命令启动 Fyne |
| 2.5 | 验证 CGO 编译 | 确保 `go build` 在开发机通过（需要 gcc） |

**验收**：
- `go run ./cmd/devmemory` 打开原生窗口
- 左侧导航可切换，右侧显示占位内容
- CLI 子命令仍然正常

---

#### Step 3：Capture 页

**目标**：Fyne Capture 页能新增 Entry，与 Web UI Capture 功能对等。

**任务**：

| # | 任务 | 说明 |
|---|------|------|
| 3.1 | 创建 `internal/ui/fyne/page_capture.go` | Capture 页布局 |
| 3.2 | Content 输入框（多行） | 主要输入区 |
| 3.3 | Type 下拉选择（Auto-detect + 11 种类型） | 可手动指定类型 |
| 3.4 | Title / Project / Tags 输入 | 可选字段 |
| 3.5 | 提交按钮 → 调用 Service.CreateEntry | 提交后清空表单 |
| 3.6 | 提交反馈（成功/失败提示） | Fyne dialog |

**验收**：
- 在 Fyne 中输入内容并提交
- `devmemory list` 能看到新增的 Entry
- 类型自动推断正确

---

#### Step 4：Search 页

**目标**：Fyne Search 页能搜索并查看 Entry。

**任务**：

| # | 任务 | 说明 |
|---|------|------|
| 4.1 | 创建 `internal/ui/fyne/page_search.go` | Search 页布局 |
| 4.2 | 搜索输入框 + 类型筛选下拉 | 回车或按钮触发搜索 |
| 4.3 | 结果列表（显示 type / title / time / score） | 按评分排序 |
| 4.4 | 点击结果 → Entry 详情弹窗 | 复用 entry_dialog.go |
| 4.5 | 创建 `internal/ui/fyne/entry_dialog.go` | Entry 详情/编辑弹窗 |

**验收**：
- 搜索关键词返回结果
- 结果按权重排序
- 点击结果可查看完整 Entry

---

#### Step 5：Today / Actions / Knowledge 页

**目标**：完成核心查看页面。

**任务**：

| # | 任务 | 说明 |
|---|------|------|
| 5.1 | 创建 `internal/ui/fyne/page_today.go` | 今日条目时间线 |
| 5.2 | 创建 `internal/ui/fyne/page_actions.go` | 可复用动作（command/url/snippet/prompt/file/folder） |
| 5.3 | 创建 `internal/ui/fyne/page_knowledge.go` | 长期知识（note/issue/business/journal/task） |
| 5.4 | Copy 按钮 → Service.CopyEntry | 复制到剪贴板 + UseCount++ |
| 5.5 | Open 按钮 → Executor.OpenURL/OpenFile | 打开 URL/文件/文件夹 |
| 5.6 | Execute 按钮（command 类型） | 确认框 → 二次确认（dangerous） → Executor.RunCommand |
| 5.7 | Entry 操作：Favorite / Archive / Edit / Delete | 统一在 entry_dialog 中 |

**验收**：
- Today 页显示当日条目
- Actions 页只显示动作类型
- Knowledge 页只显示日志类型
- Copy/Open/Execute 功能正常
- Dangerous 命令需二次确认

---

#### Step 6：Settings 页 + 导出导入

**目标**：完成设置和数据管理功能。

**任务**：

| # | 任务 | 说明 |
|---|------|------|
| 6.1 | 创建 `internal/ui/fyne/page_settings.go` | Settings 页布局 |
| 6.2 | 显示数据目录路径、条目总数、版本号 | 只读信息展示 |
| 6.3 | Export Today Markdown 按钮 | 选择保存路径 → Service.ExportTodayMarkdown |
| 6.4 | Export JSON 按钮 | 选择保存路径 → Service.ExportJSON |
| 6.5 | Import JSON 按钮 | 选择文件 → Service.ImportJSON |
| 6.6 | 状态栏（底部） | 显示版本 / 数据路径 / 条目数 |

**验收**：
- Settings 页正确显示信息
- 导出 Markdown / JSON 功能正常
- 导入 JSON 功能正常（ID 冲突覆盖）

---

#### Step 7：构建脚本 + README 更新

**目标**：Fyne 版本可构建、可分发。

**任务**：

| # | 任务 | 说明 |
|---|------|------|
| 7.1 | 更新 `scripts/build.sh` | 支持 CGO_ENABLED=1 + fyne-cross |
| 7.2 | 更新 `README.md` | 添加 Fyne UI 说明、构建依赖（gcc）、截图 |
| 7.3 | 验证四平台构建 | fyne-cross 或本机编译 |
| 7.4 | 更新 PROGRESS.md | 标记所有步骤完成 |

**验收**：
- `go build` 本机编译通过
- `fyne-cross` 交叉编译成功（Linux/macOS/Windows）
- README 包含 Fyne 相关说明

---

---

## 第三阶段：日常可用的开发者工作记忆系统

**目标**：每天工作时愿意打开它，能管理常用资源、记录排障、导出周报。

### 功能清单

| # | 功能 | 状态 | 说明 |
|---|------|------|------|
| 3.1 | Project 字段增强 | TODO | project list / filter 子命令，Fyne Project 筛选 |
| 3.2 | Issue / Business knowledge 模板 | TODO | 预设字段模板，快速填入 |
| 3.3 | Prompt / snippet 管理优化 | TODO | 语法高亮、分类、快速调用 |
| 3.4 | Command 执行确认流程 | 部分完成 | CLI 已有 dangerous 检测；Fyne 中需弹确认框 |
| 3.5 | Dangerous command 二次确认 | 部分完成 | 同上，CLI 已有 IsDangerous |
| 3.6 | Favorite 标记 + 筛选 | 部分完成 | Favorite 字段已有；筛选 UI/CLI 待做 |
| 3.7 | UseCount / LastUsedAt 统计 | **DONE** | Copy 时已维护 |
| 3.8 | Tag 管理 | TODO | autocomplete / rename / merge |
| 3.9 | Weekly report export | TODO | 扩展 export 包，按周汇总 |
| 3.10 | Portable mode 完善 | TODO | 数据迁移、多实例检测 |
| 3.11 | Entry 分页 | TODO | 大量条目时的分页加载 |
| 3.12 | 快捷键支持 | TODO | Fyne 全局快捷键（快速 Capture、快速 Search） |

### 验收标准

1. 每天使用无明显摩擦
2. 常用 command / URL / prompt 可快速复制
3. 排障过程可记录并搜索回来
4. 日报和周报可导出
5. Portable mode 可在 U 盘上运行

### 风险和砍功能边界

- Command 执行安全边界需谨慎，宁可保守
- Tag autocomplete 如果复杂则推迟，先做手动输入
- 快捷键按 Fyne 能力实现，不做超出框架能力的定制

---

## 第四阶段：稳定的跨平台个人效率与知识工具

**目标**：长期积累上千条记录、搜索仍然快、数据可迁移、三平台自用。

### 功能清单

| # | 功能 | 状态 | 说明 |
|---|------|------|------|
| 4.1 | 全文搜索优化 | TODO | Bleve 或类似纯 Go 方案，替代内存全遍历 |
| 4.2 | Entry 关联关系 | TODO | related entries，手动/自动关联 |
| 4.3 | Project timeline 视图 | TODO | 按项目维度查看条目时间线 |
| 4.4 | Markdown vault 兼容 | TODO | Obsidian 格式导入导出 |
| 4.5 | Backup / restore | TODO | 自动备份策略（定时备份、备份数量限制） |
| 4.6 | Template variables | TODO | 命令模板可填参数，如 `ping {{host}}` |

### 可选调研项

| # | 功能 | 说明 |
|---|------|------|
| O.1 | ~~桌面壳（Fyne / Wails）~~ | **已在第二阶段实现** |
| O.2 | 全局快捷键 | 系统级快捷键唤起 Capture |
| O.3 | AI 总结接口 | 可选接口，不嵌入核心流程 |
| O.4 | Workflow 自动化 | 条件触发、批量操作 |

### 验收标准

1. 1000+ 条记录搜索 < 200ms
2. 数据备份和迁移无损
3. 三平台二进制正常运行
4. 核心逻辑无 GUI 依赖

### 风险和砍功能边界

- 全文搜索如 Bleve 引入 CGO 则换成纯 Go 方案
- AI 总结仅作可选接口，不嵌入核心流程
- Entry 关联如复杂度过高则推迟

---

## 简化假设

1. 搜索用内存全遍历 + 字符串匹配，不做倒排索引（第四阶段考虑升级）
2. bbolt 单 bucket 存所有 Entry，JSON 序列化
3. Web UI 保留但不再作为主 UI 维护
4. CLI 不用 cobra，手动 flag 解析
5. 不做分页，先支持全量返回（第三阶段考虑）
6. HTTP server 纯 localhost，无 TLS，无认证
7. 默认端口 8420
8. CGO_ENABLED=1（Fyne 需要）
9. 不做签名、公证、安装包
10. 不做云同步、多用户、浏览器插件、移动端
