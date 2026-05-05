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

---

## 二、技术选型及理由

| 选择 | 理由 |
|------|------|
| **Go** | 天然跨平台交叉编译、单二进制产出、标准库足够 |
| **bbolt** | 纯 Go 实现、单文件数据库、嵌入式无需进程、事务安全 |
| **Fyne** | 纯 Go 跨平台原生 GUI 框架、无 CGO 以外的重依赖、三平台一致体验 |
| **net/http** | 标准库、Go 1.22+ 原生路由匹配、Web UI 保留备用 |
| **手动 flag 解析** | CLI 不用 cobra 等框架，手动解析参数，支持 flags 在任意位置 |

### 为什么从 Web UI 迁移到 Fyne

第一阶段使用 Local Web UI（Go embed + HTML/CSS/JS）快速验证了功能闭环。迁移到 Fyne 的原因：

1. **原生体验**：Fyne 产生真正的原生窗口，不需要浏览器
2. **离线可用**：不依赖浏览器，不占用浏览器 tab
3. **系统集成**：更好的剪贴板、文件打开、快捷键支持
4. **单进程**：不需要 HTTP server 中间层
5. **分发简单**：单二进制，双击即用

### Web UI 的处理

Web UI 代码（`internal/server/`）暂时保留，降级为备用入口。不再作为主 UI 维护。

### 为什么适合个人使用且跨平台

- 单二进制，Windows/macOS/Linux 原生窗口
- 数据全部本地，portable mode 支持 U 盘随身携带
- 无需安装、无需账号、无需网络
- Fyne 在三平台上渲染一致

---

## 三、整体架构

### 架构图

```
┌──────────────────────────────────────────────────┐
│                    main.go                        │
│            模式选择：CLI / GUI / Serve            │
├──────────┬──────────────────────┬────────────────┤
│  CLI     │     Fyne UI          │   HTTP Server   │
│  (保留)  │   (主 UI，新增)      │   (保留，降级)  │
├──────────┴──────────────────────┴────────────────┤
│                service 层（新增）                  │
│          MemoryService 统一业务接口               │
├──────────┬──────────┬──────────┬────────────────┤
│  Search  │  Action  │  Export  │    Config       │
│  Engine  │ Executor │  Engine  │    Paths        │
├──────────┴──────────┴──────────┴────────────────┤
│                   Core Model                      │
│              Entry / EntryType                    │
├──────────────────────────────────────────────────┤
│               bbolt Store                        │
└──────────────────────────────────────────────────┘
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
│       └── main.go              # 入口：模式选择（CLI / GUI / Serve）
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
│   │   ├── executor.go          # Executor 接口 + DetectType + IsDangerous
│   │   ├── executor_linux.go    # Linux: xdg-open, xclip
│   │   ├── executor_windows.go  # Windows: cmd /c start, clip
│   │   └── executor_darwin.go   # macOS: open, pbcopy
│   │
│   ├── export/
│   │   ├── markdown.go          # ExportDailyMarkdown
│   │   ├── json.go              # ExportJSON + ImportJSON
│   │   └── export_test.go       # 测试（7 cases）
│   │
│   ├── service/
│   │   └── memory_service.go    # 统一业务层，CLI 和 Fyne 共用
│   │
│   ├── server/                  # Web UI（保留，降级）
│   │   ├── server.go
│   │   ├── handlers.go
│   │   ├── exec.go
│   │   ├── exec_windows.go
│   │   └── static/
│   │       ├── index.html
│   │       ├── style.css
│   │       └── app.js
│   │
│   ├── ui/
│   │   └── fyne/                # Fyne 原生桌面 UI（新增）
│   │       ├── app.go           # Fyne app 初始化 + service 绑定
│   │       ├── main_window.go   # 主窗口 + 左侧导航
│   │       ├── page_capture.go  # Capture 页
│   │       ├── page_search.go   # Search 页
│   │       ├── page_today.go    # Today 页
│   │       ├── page_actions.go  # Actions 页
│   │       ├── page_knowledge.go# Knowledge 页
│   │       ├── page_settings.go # Settings 页
│   │       └── entry_dialog.go  # Entry 详情/编辑弹窗
│   │
│   └── config/
│       ├── config.go            # Config struct + DefaultConfig
│       ├── paths.go             # portable mode 检测 + 通用逻辑
│       ├── paths_linux.go       # ~/.config/devmemory
│       ├── paths_windows.go     # %APPDATA%/DevMemory
│       └── paths_darwin.go      # ~/Library/Application Support/DevMemory
│
├── scripts/
│   └── build.sh                 # 多平台构建（fyne-cross 或本机）
│
├── docs/
│   ├── DESIGN.md                # 架构设计（本文件）
│   ├── PLAN.md                  # 路线图与任务拆解
│   ├── SPEC.md                  # 数据结构与接口
│   └── PROGRESS.md              # 开发进度跟踪
│
└── dist/                        # 构建产物（.gitignore）
```

### 核心数据流

```
用户操作 (CLI / Fyne Widget)
    │
    ▼
┌──────────────┐
│ MemoryService│  ← CLI 和 Fyne 共用
└──────┬───────┘
       │
  ┌────┼────────────┬──────────┐
  ▼    ▼            ▼          ▼
Core  Store      Search     Action
Entry (bbolt)   Engine    Executor
```

### Entry 生命周期

```
Capture (CLI add / Fyne Capture 页 / Web POST)
  → service.CreateEntry(input)
      → action.DetectType(content)
      → core.NewEntry(type, content)
      → action.IsDangerous(content) → entry.Dangerous
      → store.Create(entry)
  → Search (service.SearchEntries → search.Engine)
  → Copy (service.CopyEntry → clipboard + UseCount++)
  → Edit (service.UpdateEntry → store.Update)
  → Archive (service.UpdateEntry → archived=true)
  → Export (service.ExportTodayMarkdown / ExportJSON)
  → Delete (service.DeleteEntry → store.Delete)
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

## 五、分层设计

### Service 层（核心新增）

Service 层是 CLI 和 Fyne UI 之间的桥梁，统一所有业务逻辑。

```go
type MemoryService struct {
    store    store.Store
    searcher *search.Engine
    executor action.Executor
}
```

**职责**：
- 封装所有业务操作（创建、查询、搜索、导出等）
- 自动类型推断、危险检测
- UseCount/LastUsedAt 维护
- 导出/导入操作

**不负责**：
- UI 展示逻辑
- 平台差异（委托给 action.Executor）
- 数据库选型（委托给 store.Store）

### CLI / Fyne / Core / Store 分层

- **Core** 定义 Entry 模型，不依赖任何存储或 UI
- **Store** 通过 interface 暴露 CRUD + ResolveID + GetByDate
- **Search** 从 Store 获取数据后自行评分排序
- **Action** 提供类型推断、危险检测、平台操作（打开/复制/执行）
- **Export** 将 Entry 列表转为 Markdown/JSON
- **Service** 统一编排 Core + Store + Search + Action + Export
- **CLI** 调用 Service，不直接操作 Store
- **Fyne UI** 调用 Service，widget 回调中不包含业务逻辑
- **Server** 调用 Service（或保持直接调用，降级维护）
- **Web UI** 纯静态文件，通过 HTTP API 交互

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
    ResolveID(prefix string) (string, error)
    GetByDate(year, month, day int) ([]*core.Entry, error)
}
```

### 平台相关逻辑隔离

```
action/executor.go           # Executor 接口 + DetectType + IsDangerous
action/executor_linux.go     # Linux: xdg-open, xclip
action/executor_windows.go   # Windows: cmd /c start, clip
action/executor_darwin.go    # macOS: open, pbcopy

config/paths.go              # 通用逻辑（portable mode 检测）
config/paths_linux.go        # ~/.config/devmemory
config/paths_windows.go      # %APPDATA%/DevMemory
config/paths_darwin.go       # ~/Library/Application Support/DevMemory
```

---

## 六、Fyne UI 设计

### 主窗口布局

```
┌─────────────────────────────────────────────┐
│  DevMemory                          [—][□][×]│
├──────────┬──────────────────────────────────┤
│ Capture  │                                  │
│ Search   │       右侧内容区域               │
│ Today    │       (各页面内容)                │
│ Actions  │                                  │
│Knowledge │                                  │
│ Settings │                                  │
├──────────┴──────────────────────────────────┤
│  状态栏：版本 / 数据路径 / 条目数           │
└─────────────────────────────────────────────┘
```

左侧导航项：

| 导航项 | 功能 |
|--------|------|
| Capture | 输入内容、选择类型/标题/项目/标签、提交 |
| Search | 搜索框 + 结果列表 + 详情预览 |
| Today | 今日条目时间线 |
| Actions | 可复用动作（command/url/snippet/prompt/file/folder） |
| Knowledge | 长期知识（note/issue/business/journal） |
| Settings | 数据目录、导出导入、版本信息 |

### 安全设计

1. command 默认只复制，不直接执行
2. 执行 command 必须弹确认框
3. `dangerous=true` 的 command 需要二次确认
4. 不保存密码、token、私钥
5. 不自动提权
6. 不静默后台危险执行

---

## 七、CLI 设计

```bash
devmemory                      # 默认启动 Fyne GUI
devmemory gui                  # 启动 Fyne GUI（显式）
devmemory serve [--port] [--no-open]  # 启动 Web UI（保留）
devmemory add <content> [--type] [--title] [--project] [--tags]
devmemory list [--type] [--project] [--tag]
devmemory search <query> [--type]
devmemory show <id>
devmemory edit <id> [--title] [--content] [--type] ...
devmemory delete <id> [--force]
devmemory today
devmemory export today|json [-o file]
devmemory import <file>
devmemory version
```

### 无参数行为

- **双击 exe / 无参数**：启动 Fyne GUI
- **CLI 子命令**：正常执行 CLI 操作
- Windows 双击行为保持向后兼容

---

## 八、数据目录策略

### 默认数据目录

| 平台 | 路径 |
|------|------|
| Linux | `~/.config/devmemory` |
| Windows | `%APPDATA%/DevMemory` |
| macOS | `~/Library/Application Support/DevMemory` |

### Portable Mode

可执行文件旁边存在 `devmemory-data/` 目录时优先使用。

---

## 九、构建目标

### 构建变化

| 项目 | 第一阶段 (Web UI) | 第二阶段 (Fyne) |
|------|-------------------|-----------------|
| CGO | CGO_ENABLED=0 | CGO_ENABLED=1 |
| 交叉编译 | go build 直接交叉 | fyne-cross 或目标平台 C 工具链 |
| 本机开发 | go build | go build（需要 gcc） |
| 依赖 | 2（bbolt + x/sys） | 2 + fyne（较多间接依赖） |

### 本机开发

```bash
# 需要 gcc
go run ./cmd/devmemory
go run ./cmd/devmemory gui
```

### 本机构建

```bash
go build -o dist/devmemory ./cmd/devmemory
```

### 交叉编译（推荐 fyne-cross）

```bash
fyne-cross windows -arch=amd64 -app-id devmemory
fyne-cross darwin -arch=amd64 -app-id devmemory
fyne-cross darwin -arch=arm64 -app-id devmemory
fyne-cross linux -arch=amd64 -app-id devmemory
```

### 限制说明

- 个人使用，不需要签名、公证
- macOS 个人使用可手动放行
- Windows 不需要 installer
- fyne-cross 需要 Docker

---

## 十、产品原则

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
native-ui
```

第一阶段（Web UI）成功的标准已达成。第二阶段（Fyne）成功的标准：

```
1. 能打开原生窗口；
2. 能新增 Entry；
3. 能搜索 Entry；
4. 能查看今日记录；
5. 能复制 snippet/prompt/command；
6. 能打开 URL/file/folder；
7. 能导出今日 Markdown；
8. 能导入/导出 JSON；
9. CLI 仍然可用；
10. 现有数据不丢失。
```
