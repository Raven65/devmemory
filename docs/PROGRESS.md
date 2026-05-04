# DevMemory — 开发进度跟踪

> 每个步骤完成后更新此文档。状态：TODO / IN PROGRESS / DONE。

---

## 当前阶段：第 1 周 — Core + Store + CLI

**开始日期**：2026-05-04
**目标**：可编译运行的 `devmemory` 二进制，支持 add / list / search / today

---

## 步骤进度

### 1. 项目骨架

| # | 任务 | 状态 |
|---|------|------|
| 1.1 | 初始化 Go module | TODO |
| 1.2 | 创建目录结构 | TODO |

### 2. Core 模型

| # | 任务 | 状态 |
|---|------|------|
| 2.1 | `internal/core/types.go` — EntryType 常量 | TODO |
| 2.2 | `internal/core/entry.go` — Entry struct + 辅助方法 | TODO |

### 3. Store 层

| # | 任务 | 状态 |
|---|------|------|
| 3.1 | `internal/store/store.go` — Store interface + ListOptions | TODO |
| 3.2 | `internal/store/bbolt_store.go` — bbolt 实现 | TODO |
| 3.3 | `internal/store/bbolt_store_test.go` — 测试 | TODO |

### 4. Config / Paths

| # | 任务 | 状态 |
|---|------|------|
| 4.1 | `internal/config/config.go` — Config struct | TODO |
| 4.2 | `internal/config/paths.go` — 通用逻辑 + portable mode | TODO |
| 4.3 | `internal/config/paths_linux.go` — Linux 路径 | TODO |
| 4.4 | `internal/config/paths_windows.go` — Windows 路径 | TODO |
| 4.5 | `internal/config/paths_darwin.go` — macOS 路径 | TODO |

### 5. Search

| # | 任务 | 状态 |
|---|------|------|
| 5.1 | `internal/search/search.go` — 搜索逻辑 | TODO |
| 5.2 | `internal/search/scorer.go` — 权重评分 | TODO |
| 5.3 | `internal/search/search_test.go` — 测试 | TODO |

### 6. CLI

| # | 任务 | 状态 |
|---|------|------|
| 6.1 | `cmd/devmemory/main.go` — 入口 + 子命令分发 | TODO |
| 6.2 | version 子命令 | TODO |
| 6.3 | add 子命令（含自动类型推断） | TODO |
| 6.4 | list 子命令 | TODO |
| 6.5 | search 子命令 | TODO |
| 6.6 | today 子命令 | TODO |

### 7. 构建与文档

| # | 任务 | 状态 |
|---|------|------|
| 7.1 | `scripts/build.sh` — 四平台构建脚本 | TODO |
| 7.2 | `README.md` — 第一版 | TODO |

---

## 集成验证

| # | 检查项 | 状态 |
|---|--------|------|
| V1 | `go build ./cmd/devmemory` 编译通过 | TODO |
| V2 | `go test ./...` 全部通过 | TODO |
| V3 | CLI add → list 可工作 | TODO |
| V4 | CLI search 可返回结果 | TODO |
| V5 | CLI today 可列出当日记录 | TODO |
| V6 | 四平台交叉编译成功 | TODO |

---

## 已知问题

（开发过程中记录遇到的 bug、设计调整、砍掉的功能等）

---

## 变更日志

| 日期 | 变更 |
|------|------|
| 2026-05-04 | 初始化项目文档（DESIGN.md, PLAN.md, SPEC.md, PROGRESS.md） |
