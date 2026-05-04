# DevMemory — 路线图与任务拆解

## 一、1 / 3 / 6 个月路线图

### 第 1 个月：最小可用个人记忆库

**目标**：能在 Linux 上日常使用，能交叉编译出三平台二进制。

**必做**：
1. Go 项目骨架 + module
2. Entry 模型 + EntryType
3. bbolt store（CRUD + 全量遍历）
4. CLI：version / add / list / search / today
5. 搜索引擎（字段匹配 + 权重排序）
6. Local HTTP server + embed 静态资源
7. Web UI：Capture / Search / Today 三个页面
8. Markdown daily export
9. JSON backup / export / import
10. 四平台构建脚本
11. Portable mode 数据目录

**不做**：
- Session / weekly report
- Command 执行（仅复制）
- 全文索引引擎
- 原生 GUI / 系统托盘 / 全局快捷键
- 浏览器插件 / 移动端
- AI 总结

**验收标准**：
1. `./devmemory add "ss -tinp" --type command --tags linux,tcp` 可新增
2. `./devmemory search "tcp"` 可搜索
3. `./devmemory today` 可列出今日记录
4. `./devmemory serve` 启动 HTTP server 并打开浏览器
5. 浏览器中可 Capture / Search / 查看今天
6. `./devmemory export today -o daily.md` 导出 Markdown
7. `./devmemory export json -o backup.json` 导出 JSON
8. `./devmemory import backup.json` 可导入
9. 四平台编译成功

**风险和砍功能边界**：
- bbolt 性能在上千条记录内无问题，暂无风险
- 如果前端 HTML/JS 工作量超预期，先只做纯文本渲染，不做复杂交互
- 搜索第一阶段用内存遍历 + 字符串匹配，不做倒排索引

---

### 第 3 个月：日常可用的开发者工作记忆系统

**目标**：每天工作时愿意打开它，能管理常用资源、记录排障、导出周报。

**必做**：
1. Project 字段增强（project list / filter）
2. Issue / Business knowledge 模板
3. Prompt / snippet 管理体验优化
4. Command 执行确认流程
5. Dangerous command 二次确认
6. Favorite 标记 + 筛选
7. UseCount / LastUsedAt 统计
8. Tag 管理（autocomplete / rename / merge）
9. Weekly report export
10. Portable mode 完善
11. Web UI 改进（Entry 编辑、分页、快捷键）

**不做**：
- 云同步 / 多用户
- 全文搜索引擎（如 Bleve）
- 原生 GUI
- 浏览器插件

**验收标准**：
1. 每天使用无明显摩擦
2. 常用 command / URL / prompt 可快速复制
3. 排障过程可记录并搜索回来
4. 日报和周报可导出
5. portable mode 可在 U 盘上运行

**风险和砍功能边界**：
- Command 执行安全边界需谨慎，宁可保守
- Tag autocomplete 如果复杂则推迟，先做手动输入
- Web UI 改进按需，不做过度设计

---

### 第 6 个月：稳定的跨平台个人效率与知识工具

**目标**：长期积累上千条记录、搜索仍然快、数据可迁移、三平台自用。

**必做**：
1. 全文搜索优化（Bleve 或类似纯 Go 方案）
2. Entry 关联关系（related entries）
3. Project timeline 视图
4. Markdown vault import/export（兼容 Obsidian 格式）
5. Backup / restore（自动备份策略）
6. Template variables（命令模板可填参数）

**可选（调研后决定）**：
- 桌面壳（Fyne / Wails）
- 全局快捷键
- AI 总结接口
- Workflow 自动化

**验收标准**：
1. 1000+ 条记录搜索 < 200ms
2. 数据备份和迁移无损
3. 三平台二进制正常运行
4. 核心逻辑无 GUI 依赖

**风险和砍功能边界**：
- 全文搜索如 Bleve 引入 CGO 则换成纯 Go 方案
- AI 总结仅作可选接口，不嵌入核心流程
- 桌面壳仅调研，不确定就推迟

---

## 二、第 1 个月详细任务拆解（按周）

### 第 1 周：Core + Store + CLI

**任务**：
1. 初始化 Go module，定义 `devmemory` 项目
2. 实现 `internal/core/entry.go` 和 `types.go` — Entry 模型
3. 实现 `internal/store/store.go` — Store interface
4. 实现 `internal/store/bbolt_store.go` — bbolt CRUD
5. 实现 `internal/config/paths.go` + 平台文件 — 数据目录
6. 实现 `internal/search/search.go` + `scorer.go` — 搜索
7. 实现 `cmd/devmemory/main.go` — CLI 入口
8. 实现 CLI 子命令：version / add / list / search / today
9. 编写 `internal/store/bbolt_store_test.go` — 基础测试
10. 编写 `scripts/build.sh`
11. 编写 `README.md` 第一版

**产出**：
- 可编译运行的 `devmemory` 二进制
- 可 add / list / search / today

**测试**：
- `go test ./internal/store/...` 通过
- `go test ./internal/search/...` 通过
- 手动测试 CLI 命令

**可运行命令**：
```bash
go build -o devmemory ./cmd/devmemory
./devmemory version
./devmemory add "ss -tinp | grep ESTAB" --type command --tags linux,tcp
./devmemory add "https://docs.kernel.org/networking/" --type url --tags linux,kernel
./devmemory list
./devmemory search "tcp"
./devmemory today
```

**风险**：
- bbolt API 熟悉成本：很低，文档清晰
- 搜索权重调优：先用简单权重，后续迭代

---

### 第 2 周：Search + Today + Markdown Export

**任务**：
1. 搜索权重优化（title > tag > project > content）
2. Entry update / delete / archive
3. CLI：show / delete / edit
4. `internal/export/markdown.go` — Markdown daily export
5. `internal/export/json.go` — JSON export / import
6. CLI：`export today` / `export json` / `import`
7. 测试：export / import round-trip

**产出**：
- 完整的 CLI 命令集
- Markdown 和 JSON 导出

**测试**：
- Export + Import round-trip 数据无损
- Markdown 输出格式正确

**可运行命令**：
```bash
./devmemory show <id>
./devmemory delete <id>
./devmemory export today -o daily.md
./devmemory export json -o backup.json
./devmemory import backup.json
```

**风险**：
- JSON import 需处理 ID 冲突（策略：保留原 ID，如已存在则覆盖）

---

### 第 3 周：Local Web UI + HTTP API

**任务**：
1. `internal/server/server.go` — HTTP server
2. `internal/server/handlers.go` — REST API handlers
3. `web/index.html` + `app.js` + `style.css` — 极简前端
4. Go embed 静态资源
5. `serve` 命令：启动 server + 自动打开浏览器
6. Capture 页面
7. Search 页面
8. Today 页面
9. Entry Detail 页面

**产出**：
- `./devmemory serve` 启动 Web UI
- 浏览器可操作所有功能

**测试**：
- curl 测试所有 API endpoint
- 浏览器手动测试 UI

**可运行命令**：
```bash
./devmemory serve
# 浏览器打开 http://127.0.0.1:8420
```

**风险**：
- 前端工作量可能超预期：保持极简，先做功能再做美观
- 浏览器自动打开在不同平台可能需要调通

---

### 第 4 周：JSON Backup + Cross Build + README

**任务**：
1. JSON import 完善和错误处理
2. Portable mode 检测
3. 四平台交叉编译脚本
4. 实际编译测试
5. README 完善（安装、使用、截图）
6. 整体 bug 修复和体验打磨
7. Entry edit（Web UI）

**产出**：
- 四个平台二进制
- 完整 README
- 可发布的个人使用版本

**测试**：
- 四平台编译成功
- Linux 实际运行测试
- Export / Import 完整流程

**可运行命令**：
```bash
./scripts/build.sh
ls -la dist/
```

**风险**：
- macOS 二进制无法在本机测试（无 macOS 环境）：编译应成功，逻辑测试靠 Linux
- Windows exe 同理

---

## 三、简化假设（第 1 个月）

1. 搜索用内存全遍历 + 字符串匹配，不做倒排索引
2. bbolt 单 bucket 存所有 Entry，JSON 序列化
3. 前端不用框架，原生 HTML/CSS/JS
4. CLI 不用 cobra，用标准库 flag + 子命令分发（或轻量 cobra）
5. 不做分页，先支持全量返回（个人使用数据量可控）
6. HTTP server 不做 TLS，纯 localhost
7. 不做用户认证（仅本地访问）
8. 默认端口 8420（可配置）

---

## 四、第 1 周交付清单

- [ ] `go.mod` 初始化
- [ ] `internal/core/types.go` — EntryType 常量
- [ ] `internal/core/entry.go` — Entry struct
- [ ] `internal/store/store.go` — Store interface
- [ ] `internal/store/bbolt_store.go` — bbolt 实现
- [ ] `internal/config/paths.go` + 平台文件 — 数据目录
- [ ] `internal/search/search.go` — 搜索逻辑
- [ ] `internal/search/scorer.go` — 权重评分
- [ ] `cmd/devmemory/main.go` — CLI 入口 + 所有子命令
- [ ] `internal/store/bbolt_store_test.go` — Store 测试
- [ ] `internal/search/search_test.go` — Search 测试
- [ ] `scripts/build.sh` — 构建脚本
- [ ] `README.md` — 第一版文档
