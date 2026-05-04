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

## 三、Search 接口

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

## 四、Config 接口

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

## 六、Export 接口

```go
// internal/export/markdown.go
func ExportDailyMarkdown(entries []*core.Entry, date time.Time) (string, error)

// internal/export/json.go
func ExportJSON(entries []*core.Entry) ([]byte, error)
func ImportJSON(data []byte) ([]*core.Entry, error)
```

---

## 七、HTTP API

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

## 八、bbolt 存储约定

- Bucket 名称：`entries`
- Key：Entry.ID (string)
- Value：JSON 序列化的 Entry
- 数据库文件：`devmemory.db`

---

## 九、CLI 子命令

| 命令 | 说明 | 状态 |
|------|------|------|
| `devmemory version` | 版本信息 | DONE |
| `devmemory serve` | 启动 HTTP server | TODO (Week 3) |
| `devmemory add` | 新增 Entry（含自动类型推断 + 危险检测） | DONE |
| `devmemory list` | 列表（支持 --type/--project/--tag 筛选） | DONE |
| `devmemory show` | 详情（支持短 ID 前缀） | DONE |
| `devmemory delete` | 删除（含确认提示，--force 跳过） | DONE |
| `devmemory edit` | 编辑（--title/--content/--type/--project/--tags/--favorite/--archive） | DONE |
| `devmemory search` | 搜索（多词 + 加权排序） | DONE |
| `devmemory today` | 今日记录（按时间线格式） | DONE |
| `devmemory export today` | 导出 Markdown（-o 输出到文件） | DONE |
| `devmemory export json` | 导出 JSON（-o 输出到文件） | DONE |
| `devmemory import` | 导入 JSON（ID 冲突时覆盖） | DONE |

---

## 十、常量与约定

| 项目 | 值 |
|------|-----|
| 默认端口 | 8420 |
| 数据库文件名 | devmemory.db |
| 导出目录 | {DataDir}/exports/ |
| 配置文件 | {DataDir}/config.json |
| Go 版本 | >= 1.23 (toolchain 1.24.11) |
| Module 名 | devmemory |
| bbolt 版本 | v1.4.3 |
| 依赖数 | 2（go.etcd.io/bbolt, golang.org/x/sys） |
