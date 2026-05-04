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
| **net/http** | 标准库、无需框架、足够本地使用 |
| **Go embed** | 将前端资源编译进二进制，实现真正的单文件分发 |

### 为什么暂时不做原生 GUI

Qt 需要 CGO，Electron 体积巨大，Wails/Fyne 增加复杂度。第一阶段的目标是"能用"，不是"好看"。Local Web UI 在三个平台上都能用浏览器打开，是最小摩擦的方案。

### 为什么适合个人使用且跨平台

- 单二进制，双击即用（Windows）或 `./devmemory` 即用（Linux/macOS）
- 数据全部本地，portable mode 支持U盘随身携带
- 无需安装、无需账号、无需网络
- Go 交叉编译从 Linux 一键产出四个平台的二进制

---

## 三、整体架构

### 模块图

```
┌─────────────────────────────────────────────────┐
│                    CLI (cobra)                   │
├─────────────────────────────────────────────────┤
│               Local Web UI (embed)              │
│         index.html / app.js / style.css         │
├─────────────────────────────────────────────────┤
│              HTTP API (net/http)                │
│  GET/POST /api/entries, /api/search, ...        │
├──────────┬──────────┬──────────┬────────────────┤
│  Search  │  Action  │  Export  │    Config      │
│  Engine  │ Executor │  Engine  │    Paths       │
├──────────┴──────────┴──────────┴────────────────┤
│                   Core Model                    │
│              Entry / EntryType                  │
├─────────────────────────────────────────────────┤
│               bbolt Store                       │
│          (data dir / portable mode)             │
└─────────────────────────────────────────────────┘
```

### 目录结构

```
devmemory/
├── go.mod
├── README.md
│
├── cmd/
│   └── devmemory/
│       └── main.go              # CLI 入口 + serve 命令
│
├── internal/
│   ├── core/
│   │   ├── entry.go             # Entry struct 定义
│   │   └── types.go             # EntryType 常量
│   │
│   ├── store/
│   │   ├── store.go             # Store interface
│   │   └── bbolt_store.go       # bbolt 实现
│   │
│   ├── search/
│   │   ├── search.go            # 搜索逻辑
│   │   └── scorer.go            # 权重评分
│   │
│   ├── action/
│   │   ├── executor.go          # 动作执行接口 + 危险检测
│   │   ├── executor_linux.go    # Linux 平台实现
│   │   ├── executor_windows.go  # Windows 平台实现
│   │   └── executor_darwin.go   # macOS 平台实现
│   │
│   ├── export/
│   │   ├── markdown.go          # Markdown daily export
│   │   └── json.go              # JSON backup/export/import
│   │
│   ├── server/
│   │   ├── server.go            # HTTP server + router
│   │   └── handlers.go          # API handlers
│   │
│   └── config/
│       ├── config.go            # 配置结构
│       ├── paths.go             # 通用路径逻辑 + portable mode
│       ├── paths_linux.go
│       ├── paths_windows.go
│       └── paths_darwin.go
│
├── web/
│   ├── index.html
│   ├── app.js
│   └── style.css
│
├── scripts/
│   └── build.sh                 # 四平台构建脚本
│
└── docs/
    ├── DESIGN.md                # 本文件
    └── PLAN.md                  # 路线图和任务拆解
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
    │ Engine  │ │Executor│ │ Engine │
    └─────────┘ └────────┘ └────────┘
```

### Entry 生命周期

```
Capture (add)
  → 存入 bbolt (CreatedAt, UpdatedAt)
  → 可选设置 Tags, Project, Favorite
    → Search (按权重排序返回)
    → Action (copy/open/execute)
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
    ID          string     // UUID
    Type        EntryType  // command / url / snippet / ...
    Title       string
    Content     string
    Summary     string
    Project     string
    Tags        []string
    Favorite    bool
    Dangerous   bool
    Archived    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
    LastUsedAt  *time.Time
    UseCount    int
}
```

### EntryType

```go
type EntryType string

const (
    EntryTypeCommand  EntryType = "command"
    EntryTypeURL      EntryType = "url"
    EntryTypeSnippet  EntryType = "snippet"
    EntryTypePrompt   EntryType = "prompt"
    EntryTypeNote     EntryType = "note"
    EntryTypeIssue    EntryType = "issue"
    EntryTypeJournal  EntryType = "journal"
    EntryTypeBusiness EntryType = "business"
    EntryTypeTask     EntryType = "task"
    EntryTypeFile     EntryType = "file"
    EntryTypeFolder   EntryType = "folder"
)
```

### 不同 Entry 类型的默认动作

| Type | 默认动作 |
|------|----------|
| command | 确认后执行，或复制命令 |
| url | 打开浏览器 |
| snippet | 复制文本 |
| prompt | 复制文本 |
| note | 查看详情 |
| issue | 查看详情 |
| journal | 查看详情 |
| business | 查看详情 |
| task | 查看详情 / 标记完成 |
| file | 打开文件 |
| folder | 打开文件夹 |

---

## 五、解耦设计

### CLI / Web UI / Core / Store 分层

- **Core** 定义 Entry 模型，不依赖任何存储或 UI
- **Store** 通过 interface 暴露 CRUD，bbolt 是实现细节
- **Search** 依赖 Core 模型和 Store interface，不依赖 HTTP 或 CLI
- **CLI** 调用 Core + Store + Search + Export，不包含业务逻辑
- **Server** 调用 Core + Store + Search + Action，暴露 HTTP API，不包含业务逻辑
- **Web UI** 纯静态文件，只通过 HTTP API 交互

### Store Interface

```go
type Store interface {
    Create(entry *core.Entry) error
    Get(id string) (*core.Entry, error)
    Update(entry *core.Entry) error
    Delete(id string) error
    List(opts ListOptions) ([]*core.Entry, error)
    Search(query string) ([]*core.Entry, error)
}
```

### 平台相关逻辑隔离

```
executor.go         → 接口定义 + 危险命令检测（通用逻辑）
executor_linux.go   → xdg-open, exec.Command("sh", "-c", ...)
executor_windows.go → rundll32 / start, exec.Command("cmd", "/c", ...)
executor_darwin.go  → open, exec.Command("open", ...)

paths.go            → 通用逻辑（portable mode 检测）
paths_linux.go      → ~/.config/devmemory
paths_windows.go    → %APPDATA%/DevMemory
paths_darwin.go     → ~/Library/Application Support/DevMemory
```

Go 编译时通过文件名后缀自动选择对应平台文件。

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

GET    /api/entries               # 列表（支持 ?type=&project=&tag= 筛选）
POST   /api/entries               # 新增
GET    /api/entries/:id           # 详情
PUT    /api/entries/:id           # 更新
DELETE /api/entries/:id           # 删除

GET    /api/search?q=             # 搜索
GET    /api/today                 # 今日记录

POST   /api/entries/:id/execute   # 执行命令（带确认）
POST   /api/entries/:id/copy      # 记录使用次数
POST   /api/entries/:id/open      # 打开 URL/文件/文件夹

GET    /api/export/today          # 导出今日 Markdown
GET    /api/export/json           # 导出全量 JSON
POST   /api/import/json           # 导入 JSON
```

---

## 八、CLI 设计

```bash
devmemory serve                                            # 启动 HTTP server
devmemory add <content> [--type] [--title] [--project] [--tags]  # 新增
devmemory list [--type] [--project] [--tag]                # 列表
devmemory search <query>                                   # 搜索
devmemory show <entry-id>                                  # 详情
devmemory delete <entry-id>                                # 删除
devmemory today                                            # 今日记录
devmemory export today -o daily.md                         # 导出今日 Markdown
devmemory export json -o backup.json                       # 导出 JSON
devmemory import backup.json                               # 导入 JSON
devmemory version                                          # 版本信息
```

### 自动类型推断规则

```
如果内容以 http:// 或 https:// 开头 → type=url
如果内容看起来像 shell 命令（含 |, >, <, &&, ||, ;, $ 等）→ type=command
否则 → type=note
```

---

## 九、安全设计

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

## 十、构建目标

```bash
# Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/devmemory-linux-amd64 ./cmd/devmemory

# Windows amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/devmemory-windows-amd64.exe ./cmd/devmemory

# macOS Intel
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/devmemory-darwin-amd64 ./cmd/devmemory

# macOS Apple Silicon
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/devmemory-darwin-arm64 ./cmd/devmemory
```

---

## 十一、产品原则

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
