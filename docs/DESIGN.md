# DevMemory — 架构设计文档

## 一、项目定位

### 解决什么问题

程序员的日常被碎片信息淹没：一条 shell 命令、一个调试 URL、一段 prompt、一次排障过程、一条业务知识。这些信息散落在聊天记录、浏览器书签、临时 sticky note、甚至脑子里，几天后几乎不可能找回。

DevMemory 要做的是：**一个本地运行的、零摩擦的信息捕获器 + 搜索器 + 导出器**。

核心闭环：

```
随手记录 → 自动/手动分类 → 快速搜索 → 复制/打开/执行 → 日报/周报导出
```

### 不解决什么

- 云同步、多用户、团队协作
- 密码管理
- 完整的笔记系统 / Notion 替代品
- 浏览器插件 / 移动端 / 系统托盘 / 全局快捷键
- 公开发布、安装包、签名公证
- AI 自动总结（至少第一阶段）

### QuickPilot 和 DevFlow Journal 如何融合

用一个统一模型 **Entry** 承载所有类型。QuickPilot 的动作库（command / snippet / prompt / URL / file / folder）和 DevFlow Journal 的日志（note / issue / journal / business / task）都是 Entry 的不同 Type。它们共享同一套搜索、标签、项目管理能力。区别只在于：

- 动作类 Entry 有"执行"语义（复制、打开、运行）
- 日志类 Entry 有"查看"语义（阅读、导出）

这避免了过早拆分实体，同时保持模型简洁。

---

## 二、技术选型及理由

| 选择 | 理由 |
|------|------|
| **Go** | 天然跨平台交叉编译、单二进制产出、CGO_ENABLED=0 编译零依赖、标准库足够 |
| **bbolt** | 纯 Go 实现、单文件数据库、嵌入式无需进程、事务安全、key-value 足够 Entry 存储 |
| **Local Web UI** | HTML/CSS/JS 是最轻量的跨平台 UI 方案、Go embed 嵌入后零外部依赖、浏览器即界面 |
| **net/http** | 标准库、Go 1.22+ 原生路由匹配、无需框架 |
| **Go embed** | 将前端资源编译进二进制，实现真正的单文件分发 |
| **手动 flag 解析** | CLI 不用 cobra 等框架，手动解析参数，支持 flags 在任意位置，更轻量 |

### 为什么暂时不做原生 GUI

Qt 需要 CGO，Electron 体积巨大，Wails/Fyne 增加复杂度。第一阶段的目标是"能用"，不是"好看"。Local Web UI 在三个平台上都能用浏览器打开，是最小摩擦的方案。

### 为什么适合个人使用且跨平台

- 单二进制，Windows 双击自动启动 Web UI
- Linux/macOS 命令行直接运行
- 数据全部本地，portable mode 支持 U 盘随身携带
- 无需安装、无需账号、无需网络
- Go 交叉编译从 Linux 一键产出四个平台的二进制

---

## 三、整体架构

### 模块图

```
┌─────────────────────────────────────────────────┐
│                   CLI (main.go)                  │
│   add / list / show / delete / edit / search    │
│   today / export / import / serve               │
├─────────────────────────────────────────────────┤
│               Local Web UI (embed)              │
│         index.html / app.js / style.css         │
├─────────────────────────────────────────────────┤
│              HTTP API (net/http)                │
│  GET/POST /api/entries, /api/search, ...        │
├──────────┬──────────┬──────────┬────────────────┤
│  Search  │  Action  │  Export  │    Config      │
│  Engine  │ Detector │  Engine  │    Paths       │
├──────────┴──────────┴──────────┴────────────────┤
│                   Core Model                    │
│              Entry / EntryType                  │
├─────────────────────────────────────────────────┤
│               bbolt Store                       │
│      CRUD + ResolveID + GetByDate              │
└─────────────────────────────────────────────────┘
```

### 目录结构

```
devmemory/
├── go.mod
├── go.sum
├── .gitignore
├── README.md
│
├── cmd/
│   └── devmemory/
│       └── main.go              # CLI 入口 + 全部子命令 + serve
│
├── internal/
│   ├── core/
│   │   ├── entry.go             # Entry struct + NewEntry + TitleOrContent
│   │   └── types.go             # EntryType 常量（11 种）
│   │
│   ├── store/
│   │   ├── store.go             # Store interface + ListOptions
│   │   ├── bbolt_store.go       # bbolt 实现 + ResolveID + GetByDate
│   │   └── bbolt_store_test.go  # 测试（11 cases）
│   │
│   ├── search/
│   │   ├── search.go            # 搜索逻辑 + Result 结构
│   │   ├── scorer.go            # 权重评分
│   │   └── search_test.go       # 测试（7 cases）
│   │
│   ├── action/
│   │   └── executor.go          # 类型推断 DetectType + 危险检测 IsDangerous
│   │
│   ├── export/
│   │   ├── markdown.go          # ExportDailyMarkdown
│   │   ├── json.go              # ExportJSON + ImportJSON
│   │   └── export_test.go       # 测试（7 cases）
│   │
│   ├── server/
│   │   ├── server.go            # HTTP server + embed + 路由 + 优雅关闭
│   │   ├── handlers.go          # 全部 API handler
│   │   ├── exec.go              # 浏览器打开（!windows）
│   │   ├── exec_windows.go      # 浏览器打开（windows）
│   │   └── static/              # go:embed 嵌入
│   │       ├── index.html       # Web UI 页面（三视图 + 弹窗）
│   │       ├── style.css        # 暗色主题（GitHub 风格）
│   │       └── app.js           # 前端交互逻辑
│   │
│   └── config/
│       ├── config.go            # Config struct + DefaultConfig
│       ├── paths.go             # portable mode 检测 + 通用逻辑
│       ├── paths_linux.go       # ~/.config/devmemory
│       ├── paths_windows.go     # %APPDATA%/DevMemory
│       └── paths_darwin.go      # ~/Library/Application Support/DevMemory
│
├── scripts/
│   └── build.sh                 # 四平台构建 → dist/
│
├── docs/
│   ├── DESIGN.md                # 架构设计（本文件）
│   ├── PLAN.md                  # 路线图与任务拆解
│   ├── SPEC.md                  # 数据结构与接口
│   └── PROGRESS.md              # 开发进度跟踪
│
└── dist/                        # 构建产物（.gitignore）
    ├── devmemory-linux-amd64
    ├── devmemory-windows-amd64.exe
    ├── devmemory-darwin-amd64
    └── devmemory-darwin-arm64
```

### 核心数据流

```
用户输入 (CLI / Web UI)
    │
    ▼
┌──────────┐    ┌─────────┐    ┌──────────┐
│ CLI flag │───▶│  Core   │───▶│  Store   │
│ 解析     │    │ Entry   │    │ (bbolt)  │
└──────────┘    └─────────┘    └──────────┘
                    │
         ┌──────────┼──────────┐
         ▼          ▼          ▼
    ┌─────────┐ ┌────────┐ ┌────────┐
    │ Search  │ │ Action │ │ Export │
    │ Engine  │ │Detect  │ │ Engine │
    └─────────┘ └────────┘ └────────┘
```

### Entry 生命周期

```
Capture (add / Web UI POST)
  → 自动类型推断 + 危险检测
  → 存入 bbolt (CreatedAt, UpdatedAt)
  → 可选设置 Tags, Project, Favorite
    → Search (按权重排序返回)
    → Copy (UseCount++, LastUsedAt 更新)
    → Edit (修改内容, UpdatedAt 更新)
    → Archive (Archived=true)
    → Export (Markdown / JSON)
    → Delete (物理删除)
```

---

## 四、核心数据模型

### Entry

```go
type Entry struct {
    ID          string     // 32 字符 hex ID (crypto/rand)
    Type        EntryType  // command / url / snippet / ...
    Title       string
    Content     string
    Summary     string
    Project     string
    Tags        []string
    Favorite    bool
    Dangerous   bool       // 自动检测：rm -rf, del /s, format, ...
    Archived    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
    LastUsedAt  *time.Time // Copy 时更新
    UseCount    int        // Copy 时递增
}
```

### EntryType（11 种）

| Type | 分类 | 默认动作 | 说明 |
|------|------|----------|------|
| command | 动作 | 确认后执行，或复制 | shell 命令 |
| url | 动作 | 打开浏览器 | 网址 |
| snippet | 动作 | 复制文本 | 代码片段 |
| prompt | 动作 | 复制文本 | AI prompt |
| file | 动作 | 打开文件 | 文件路径 |
| folder | 动作 | 打开文件夹 | 文件夹路径 |
| note | 日志 | 查看详情 | 随手笔记 |
| issue | 日志 | 查看详情 | 问题记录 |
| journal | 日志 | 查看详情 | 日志条目 |
| business | 日志 | 查看详情 | 业务知识 |
| task | 日志 | 查看详情 / 标记完成 | 待办任务 |

---

## 五、解耦设计

### CLI / Web UI / Core / Store 分层

- **Core** 定义 Entry 模型，不依赖任何存储或 UI
- **Store** 通过 interface 暴露 CRUD + ResolveID + GetByDate，bbolt 是实现细节
- **Search** 依赖 Core 模型，从 Store 获取数据后自行评分排序，不依赖 HTTP 或 CLI
- **Action** 依赖 Core 模型，提供类型推断和危险检测，不依赖存储
- **Export** 依赖 Core 模型，将 Entry 列表转为 Markdown/JSON，不依赖存储
- **CLI** 调用 Core + Store + Search + Export + Server，不包含业务逻辑
- **Server** 调用 Core + Store + Search + Export，暴露 HTTP API，embed 静态文件
- **Web UI** 纯静态文件（HTML/CSS/JS），只通过 HTTP API 交互

### Store Interface

```go
type Store interface {
    Open() error
    Close() error
    Create(entry *core.Entry) error
    Get(id string) (*core.Entry, error)
    Update(entry *core.Entry) error
    Delete(id string) error
    List(opts ListOptions) ([]*core.Entry, error)
    ResolveID(prefix string) (string, error)    // 短 ID → 完整 ID
    GetByDate(year, month, day int) ([]*core.Entry, error)
}
```

### 平台相关逻辑隔离

```
server/exec.go           → 浏览器打开（!windows，使用 xdg-open/open）
server/exec_windows.go   → 浏览器打开（windows，使用 cmd /c start）

config/paths.go          → 通用逻辑（portable mode 检测）
config/paths_linux.go    → ~/.config/devmemory
config/paths_windows.go  → %APPDATA%/DevMemory
config/paths_darwin.go   → ~/Library/Application Support/DevMemory
```

Go 编译时通过 build tag（`//go:build windows`）自动选择对应平台文件。

---

## 六、数据目录策略

### 默认数据目录

| 平台 | 路径 |
|------|------|
| Linux | `~/.config/devmemory` |
| Windows | `%APPDATA%/DevMemory` |
| macOS | `~/Library/Application Support/DevMemory` |

### Portable Mode

如果可执行文件旁边存在 `devmemory-data/` 目录，则优先使用它。

Portable 目录结构：

```
devmemory.exe
devmemory-data/
  devmemory.db
  exports/
  config.json
```

---

## 七、HTTP API 设计

```
GET    /api/health                # 健康检查
GET    /                          # Web UI（embed 静态文件）

GET    /api/entries               # 列表（?type=&project=&tag= 筛选）
POST   /api/entries               # 新增（JSON body）
GET    /api/entries/{id}          # 详情（支持短 ID 前缀）
PUT    /api/entries/{id}          # 更新（partial update）
DELETE /api/entries/{id}          # 删除（支持短 ID 前缀）

GET    /api/search?q=&type=       # 搜索（返回 entry + score）
GET    /api/today                 # 今日记录（date + entries）

POST   /api/entries/{id}/copy     # 记录使用 + 返回 content

GET    /api/export/today          # 导出今日 Markdown（attachment）
GET    /api/export/json           # 导出全量 JSON（attachment）
POST   /api/import/json           # 导入 JSON（数组或 {entries:[...]})
```

所有 API 端点加 CORS 头（`Access-Control-Allow-origin: *`），便于本地开发。

---

## 八、CLI 设计

```bash
devmemory serve [--port 8420] [--no-open]                   # 启动 Web UI
devmemory add <content> [--type] [--title] [--project] [--tags]
devmemory list [--type] [--project] [--tag]                  # 别名：ls
devmemory search <query> [--type]                            # 别名：s
devmemory show <id>                                          # 支持短 ID
devmemory edit <id> [--title] [--content] [--type] ...       # 支持短 ID
devmemory delete <id> [--force]                              # 别名：rm，支持短 ID
devmemory today                                              # 今日记录
devmemory export today [-o file]                             # Markdown
devmemory export json [-o file]                              # JSON
devmemory import <file>                                      # JSON（ID 冲突覆盖）
devmemory version                                            # 版本信息
```

### Windows 双击行为

Windows 上双击 `devmemory.exe`（无参数）自动执行 `serve` 命令，启动 Web UI 并打开浏览器。其他平台无参数时显示帮助信息。

### 自动类型推断规则

```
如果内容以 http:// 或 https:// 开头 → type=url
如果内容看起来像 shell 命令（含 |, >, <, &&, ||, ;, $ 等）→ type=command
否则 → type=note
```

---

## 九、Web UI 设计

三个视图 + 一个弹窗：

| 视图 | 功能 |
|------|------|
| **Capture** | 输入内容、选择类型/标题/项目/标签、提交 |
| **Search** | 搜索框 + 类型筛选、结果列表（含评分） |
| **Today** | 今日条目时间线、刷新按钮 |

**Entry 详情弹窗**：点击任意条目卡片弹出，显示完整信息，支持 Copy / Favorite / Archive / Delete 操作。

暗色主题，GitHub 风格配色（`#0d1117` 背景，`#58a6ff` 强调色），响应式布局。

---

## 十、安全设计

1. command 默认不直接执行，必须确认
2. `dangerous=true` 的 command 需要二次确认
3. 不保存密码、token、私钥
4. 不自动提权
5. 不做静默后台危险执行
6. 对明显危险命令给出警告

### 危险命令关键词

```
rm -rf
del /s
format
shutdown
mkfs
:(){ :|:& };:
dd if=
> /dev/sd
chmod -R 777 /
```

---

## 十一、构建目标

```bash
# 使用构建脚本（输出到 dist/）
./scripts/build.sh [version]

# 手动构建
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/devmemory-linux-amd64 ./cmd/devmemory
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/devmemory-windows-amd64.exe ./cmd/devmemory
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/devmemory-darwin-amd64 ./cmd/devmemory
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/devmemory-darwin-arm64 ./cmd/devmemory
```

---

## 十二、产品原则

```
Capture first, organize later.
```

```
local-first
personal-use
capture-first
search-first
cross-platform
low-dependency
single-binary-friendly
```

第一阶段成功的标准不是功能多，而是：

```
我学到一个命令，可以 5 秒内记进去；
我遇到一个问题，可以随手记录排查过程；
我几天后能搜回来；
我能导出今天做了什么；
Linux/Windows/macOS 都能跑。
```
