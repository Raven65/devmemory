# DevMemory — 数据结构与接口文档

> 随开发持续更新，记录所有核心类型、接口定义和关键约定。

---

## 一、核心类型

### EntryType

```go
// internal/core/types.go
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

### Entry

```go
// internal/core/entry.go
type Entry struct {
    ID         string     `json:"id"`
    Type       EntryType  `json:"type"`
    Title      string     `json:"title"`
    Content    string     `json:"content"`
    Summary    string     `json:"summary,omitempty"`
    Project    string     `json:"project,omitempty"`
    Tags       []string   `json:"tags,omitempty"`
    Favorite   bool       `json:"favorite"`
    Dangerous  bool       `json:"dangerous"`
    Archived   bool       `json:"archived"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    LastUsedAt *time.Time `json:"last_used_at,omitempty"`
    UseCount   int        `json:"use_count"`
}
```

### ListOptions

```go
// internal/store/store.go
type ListOptions struct {
    Type     EntryType
    Project  string
    Tag      string
    Archived bool
    Limit    int
    Offset   int
}
```

---

## 二、Store 接口

```go
// internal/store/store.go
type Store interface {
    // CRUD
    Create(entry *core.Entry) error
    Get(id string) (*core.Entry, error)
    Update(entry *core.Entry) error
    Delete(id string) error

    // 查询
    List(opts ListOptions) ([]*core.Entry, error)
    GetByDate(year, month, day int) ([]*core.Entry, error)

    // ID 前缀匹配
    ResolveID(prefix string) (string, error)

    // 生命周期
    Open() error
    Close() error
}
```

---

## 三、Service 层（新增）

Service 层是 CLI 和 Fyne UI 之间的桥梁，统一所有业务逻辑。

### MemoryService

```go
// internal/service/memory_service.go
type MemoryService struct {
    store    store.Store
    searcher *search.Engine
    executor action.Executor
}

func NewMemoryService(store store.Store, searcher *search.Engine, executor action.Executor) *MemoryService
```

### Service 方法

```go
// 创建
func (s *MemoryService) CreateEntry(input CreateEntryInput) (*core.Entry, error)

type CreateEntryInput struct {
    Content string
    Type    core.EntryType // 空则自动推断
    Title   string
    Project string
    Tags    []string
}

// 查询
func (s *MemoryService) GetEntry(id string) (*core.Entry, error)        // 支持 ResolveID 短前缀
func (s *MemoryService) ListEntries(opts store.ListOptions) ([]*core.Entry, error)
func (s *MemoryService) SearchEntries(query string, entryType core.EntryType) ([]*search.SearchResult, error)
func (s *MemoryService) GetToday() ([]*core.Entry, error)

// 操作
func (s *MemoryService) CopyEntry(id string) (*core.Entry, error)       // 复制到剪贴板 + UseCount++
func (s *MemoryService) OpenEntry(id string) error                       // 打开 URL/文件/文件夹
func (s *MemoryService) ExecuteEntry(id string) error                    // 执行命令（含确认流程）

// 修改
func (s *MemoryService) UpdateEntry(id string, fields UpdateEntryFields) (*core.Entry, error)
func (s *MemoryService) ToggleFavorite(id string) (*core.Entry, error)
func (s *MemoryService) ToggleArchive(id string) (*core.Entry, error)
func (s *MemoryService) DeleteEntry(id string) error

// 导出
func (s *MemoryService) ExportTodayMarkdown() (string, error)
func (s *MemoryService) ExportJSON() ([]byte, error)
func (s *MemoryService) ImportJSON(data []byte) (int, error)            // 返回导入条目数
```

### 职责边界

| 层 | 负责 | 不负责 |
|----|------|--------|
| **Service** | 业务编排、类型推断、危险检测、UseCount 维护 | UI 展示、平台差异 |
| **Store** | 数据持久化、CRUD、ID 解析 | 业务逻辑 |
| **Search** | 多词匹配、权重评分、结果排序 | 数据存储 |
| **Action** | 平台操作（打开/复制/执行）、类型推断、危险检测 | 业务流程 |
| **UI (Fyne/CLI)** | 用户交互、输入收集、结果展示 | 业务逻辑 |

---

## 四、Search 接口

```go
// internal/search/search.go
type Engine struct {
    store store.Store
}

type SearchResult struct {
    Entry  *core.Entry
    Score  float64
    Fields []string // 命中的字段列表
}

func (e *Engine) Search(query string, opts store.ListOptions) ([]*SearchResult, error)
```

### 搜索权重

| 字段 | 权重 |
|------|------|
| Title | 10.0 |
| Tag | 7.0 |
| Project | 5.0 |
| Summary | 3.0 |
| Content | 1.0 |
| Favorite 加成 | +2.0 |
| UseCount 加成 | +0.1 × count |
| LastUsedAt 加成 | 越近越高，上限 +3.0 |

---

## 五、Action 接口

```go
// internal/action/executor.go
type Executor interface {
    OpenURL(url string) error
    OpenFile(path string) error
    OpenFolder(path string) error
    RunCommand(cmd string) error
    CopyToClipboard(text string) error
}

func NewExecutor() Executor
func IsDangerous(cmd string) bool
func DetectType(content string) core.EntryType
```

### 平台实现

| 文件 | 平台 | 实现 |
|------|------|------|
| `executor_linux.go` | Linux | xdg-open, xclip |
| `executor_windows.go` | Windows | cmd /c start, clip |
| `executor_darwin.go` | macOS | open, pbcopy |

### 危险命令关键词

```
rm -rf, del /s, format, shutdown, mkfs, dd if=, > /dev/sd, chmod -R 777 /
```

### 自动类型推断

```
http:// 或 https:// 开头 → url
含 |, >, <, &&, ||, ;, $ 等符号 → command
否则 → note
```

---

## 六、Config 接口

```go
// internal/config/config.go
type Config struct {
    DataDir string `json:"data_dir"`
    Port    int    `json:"port"`
    DBPath  string `json:"-"`
}

// paths.go + paths_{linux,windows,darwin}.go
func DefaultConfig() *Config
func DataDir() string
func IsPortableMode() bool
```

### 平台数据目录

| 平台 | 默认路径 |
|------|----------|
| Linux | `~/.config/devmemory` |
| Windows | `%APPDATA%\DevMemory` |
| macOS | `~/Library/Application Support/DevMemory` |

Portable mode：可执行文件同目录下存在 `devmemory-data/` 时优先使用。

---

## 七、Export 接口

```go
// internal/export/markdown.go
func ExportDailyMarkdown(entries []*core.Entry, date time.Time) (string, error)

// internal/export/json.go
func ExportJSON(entries []*core.Entry) ([]byte, error)
func ImportJSON(data []byte) ([]*core.Entry, error)
```

---

## 八、Fyne UI 页面规格

### 主窗口

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

### 页面规格

| 页面 | 文件 | 功能 |
|------|------|------|
| Capture | `page_capture.go` | Content 多行输入 + Type 下拉 + Title/Project/Tags + 提交 |
| Search | `page_search.go` | 搜索框 + 类型筛选 + 结果列表 + 详情弹窗 |
| Today | `page_today.go` | 今日条目时间线 + 刷新 |
| Actions | `page_actions.go` | 动作类型列表（command/url/snippet/prompt/file/folder）+ Copy/Open/Execute |
| Knowledge | `page_knowledge.go` | 知识类型列表（note/issue/business/journal/task）+ 查看/编辑 |
| Settings | `page_settings.go` | 数据目录 + 条目数 + 版本 + 导出导入按钮 |

### Entry 弹窗

| 文件 | 功能 |
|------|------|
| `entry_dialog.go` | Entry 详情查看 + 编辑模式 + Favorite / Archive / Delete 操作 |

### 安全设计

1. command 默认只复制，不直接执行
2. 执行 command 必须弹确认框
3. `dangerous=true` 的 command 需要二次确认
4. 不保存密码、token、私钥
5. 不自动提权
6. 不静默后台危险执行

---

## 九、HTTP API（保留，降级）

| Method | Path | 说明 |
|--------|------|------|
| GET | /api/health | 健康检查 |
| GET | /api/entries | 列表（?type=&project=&tag=） |
| POST | /api/entries | 新增 |
| GET | /api/entries/:id | 详情 |
| PUT | /api/entries/:id | 更新 |
| DELETE | /api/entries/:id | 删除 |
| GET | /api/search?q= | 搜索 |
| GET | /api/today | 今日记录 |
| POST | /api/entries/:id/execute | 执行命令 |
| POST | /api/entries/:id/copy | 记录使用 |
| POST | /api/entries/:id/open | 打开 |
| GET | /api/export/today | 导出今日 Markdown |
| GET | /api/export/json | 导出全量 JSON |
| POST | /api/import/json | 导入 JSON |

---

## 十、CLI 子命令

| 命令 | 说明 | 状态 |
|------|------|------|
| `devmemory` | 默认启动 Fyne GUI | TODO |
| `devmemory gui` | 启动 Fyne GUI（显式） | TODO |
| `devmemory serve [--port] [--no-open]` | 启动 Web UI | DONE |
| `devmemory add <content> [--type] [--title] [--project] [--tags]` | 新增 Entry | DONE |
| `devmemory list [--type] [--project] [--tag]` | 列表 | DONE |
| `devmemory search <query> [--type]` | 搜索 | DONE |
| `devmemory show <id>` | 详情（短 ID 前缀） | DONE |
| `devmemory edit <id> [--title] [--content] [--type] ...` | 编辑 | DONE |
| `devmemory delete <id> [--force]` | 删除 | DONE |
| `devmemory today` | 今日记录 | DONE |
| `devmemory export today\|json [-o file]` | 导出 | DONE |
| `devmemory import <file>` | 导入 | DONE |
| `devmemory version` | 版本信息 | DONE |

---

## 十一、bbolt 存储约定

- Bucket 名称：`entries`
- Key：Entry.ID (string)
- Value：JSON 序列化的 Entry
- 数据库文件：`devmemory.db`

---

## 十二、常量与约定

| 项目 | 值 |
|------|-----|
| 默认端口 | 8420 |
| 数据库文件名 | devmemory.db |
| 导出目录 | {DataDir}/exports/ |
| 配置文件 | {DataDir}/config.json |
| Go 版本 | >= 1.23 (toolchain 1.24.11) |
| CGO | CGO_ENABLED=1（Fyne 需要） |
| Module 名 | devmemory |
| bbolt 版本 | v1.4.3 |
| 当前依赖 | go.etcd.io/bbolt, golang.org/x/sys |
| 新增依赖 | fyne.io/fyne/v2（迁移 Step 2） |
