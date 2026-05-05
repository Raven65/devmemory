# DevMemory — Windows 手动测试用例

> 第二阶段（Fyne 迁移）完成后，基于 Windows 产物的手动测试清单。

## 构建说明

在 Windows 上构建需要 GCC（推荐 [MSYS2](https://www.msys2.org/) 或 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)）：

```cmd
# 确保 gcc 在 PATH 中
gcc --version

# 构建
go build -o devmemory.exe ./cmd/devmemory
```

或使用 fyne-cross（需 Docker）：

```cmd
fyne-cross --targets=windows/amd64 .
```

---

## 一、CLI 基础功能

### 1.1 版本信息
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 1 | 双击 `devmemory.exe` | 启动 Fyne GUI 窗口（不再闪退） |
| 2 | 打开 cmd，运行 `devmemory.exe version` | 输出版本号、commit hash、构建日期 |

### 1.2 添加条目
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 3 | `devmemory.exe add "echo hello" --type command --tags test,windows` | 成功，输出新条目 ID 和类型 |
| 4 | `devmemory.exe add "https://learn.microsoft.com" --type url --tags docs,winapi` | 成功，类型 url |
| 5 | `devmemory.exe add "Windows registry keys affect service behavior" --type note --project win32` | 成功，类型 note |
| 6 | `devmemory.exe add "dir /s /b *.go"` | 自动检测为 command（含 `*` 模式） |
| 7 | `devmemory.exe add "https://go.dev/doc/install"` | 自动检测为 url（以 http 开头） |
| 8 | `devmemory.exe add "some random thought"` | 自动检测为 note（默认类型） |

### 1.3 列出条目
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 9 | `devmemory.exe list` | 显示所有条目，含 ID/类型/标题/时间 |
| 10 | `devmemory.exe list --type command` | 仅显示 command 类型 |
| 11 | `devmemory.exe list --tag test` | 仅显示含 test 标签的条目 |
| 12 | `devmemory.exe list --project win32` | 仅显示 win32 项目的条目 |

### 1.4 搜索
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 13 | `devmemory.exe search "hello"` | 找到步骤 3 添加的命令 |
| 14 | `devmemory.exe search "microsoft"` | 找到步骤 4 添加的 URL |
| 15 | `devmemory.exe search "registry"` | 找到步骤 5 添加的笔记 |
| 16 | `devmemory.exe search "nonexistent_xyz"` | 返回空结果 |

### 1.5 查看详情
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 17 | `devmemory.exe show <id>` | 显示条目完整信息（ID/类型/标题/内容/项目/标签/时间等） |
| 18 | `devmemory.exe show <短ID前缀>` | 使用 ID 前几个字符也能匹配到 |

### 1.6 编辑条目
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 19 | `devmemory.exe edit <id> --title "Updated Title"` | 成功，title 更新 |
| 20 | `devmemory.exe edit <id> --tags "new,tags"` | 成功，tags 更新 |
| 21 | `devmemory.exe edit <id> --project newproject` | 成功，project 更新 |

### 1.7 删除条目
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 22 | `devmemory.exe delete <id>` | 显示确认提示，输入 y 确认删除 |
| 23 | `devmemory.exe delete <id> --force` | 跳过确认直接删除 |

### 1.8 Today 命令
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 24 | `devmemory.exe today` | 列出今天添加的所有条目 |

### 1.9 导出/导入
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 25 | `devmemory.exe export today -o today.md` | 生成 today.md，内容为 Markdown 格式的今日条目 |
| 26 | 打开 `today.md` 检查内容 | 格式正确，条目内容完整 |
| 27 | `devmemory.exe export json -o backup.json` | 生成 backup.json |
| 28 | 打开 `backup.json` 检查内容 | JSON 数组，每个条目含所有字段 |
| 29 | `devmemory.exe import backup.json` | 成功导入，ID 相同时覆盖 |
| 30 | `devmemory.exe list` | 条目数量与导入前一致 |

---

## 二、Fyne GUI 功能

### 2.1 窗口启动
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 31 | 双击 `devmemory.exe` | 打开 Fyne GUI 窗口 |
| 32 | 检查窗口标题 | 标题为 "DevMemory" |
| 33 | 检查左侧导航 | 显示 6 个页面：Capture / Search / Today / Actions / Knowledge / Settings |
| 34 | 检查底部状态栏 | 显示版本号、数据路径、条目数 |

### 2.2 Capture 页
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 35 | 点击左侧 "Capture" | 显示录入表单 |
| 36 | 在 Content 输入 `netstat -an | findstr LISTEN` | 多行输入框可用 |
| 37 | Type 选择 `command` | 下拉选择正常 |
| 38 | 在 Title 输入 "Check listening ports" | 输入正常 |
| 39 | 在 Project 输入 "network" | 输入正常 |
| 40 | 在 Tags 输入 "windows,network,debug" | 输入正常 |
| 41 | 点击 "Save" 按钮 | 成功提示，表单清空 |
| 42 | 再添加一条：Content=`https://devblogs.microsoft.com`，Type=`url` | 成功 |
| 43 | 添加一条：Content=`PowerShell pipeline variable behavior`，Type=`note`，Project=`ps` | 成功 |

### 2.3 Search 页
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 44 | 切换到 "Search" 页 | 显示搜索框 + 类型筛选 + 空结果区 |
| 45 | 在搜索框输入 "netstat" 后回车 | 显示匹配结果 |
| 46 | 点击搜索结果中的一条 | 弹出 Entry 详情对话框 |
| 47 | 在搜索框输入 "microsoft" | 显示 URL 条目结果 |
| 48 | 使用 Type 下拉筛选为 `command` | 仅显示 command 类型的结果 |
| 49 | 搜索不存在的关键词 "zzz_nothing" | 显示空结果 |

### 2.4 Today 页
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 50 | 切换到 "Today" 页 | 显示今天添加的所有条目（包括 CLI 和 GUI 添加的） |
| 51 | 点击某一条目 | 弹出 Entry 详情对话框 |
| 52 | 点击 "Refresh" 按钮 | 列表刷新 |

### 2.5 Actions 页
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 53 | 切换到 "Actions" 页 | 显示所有 action 类型条目（command/url/snippet/prompt/file/folder） |
| 54 | 检查列表内容 | 之前 CLI 和 GUI 添加的 command、url 条目都在 |
| 55 | 点击一条 command 条目 | 弹出详情对话框 |

### 2.6 Knowledge 页
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 56 | 切换到 "Knowledge" 页 | 显示所有 knowledge 类型条目（note/issue/business/journal/task） |
| 57 | 检查列表内容 | 之前添加的 note 条目都在 |
| 58 | 点击一条 note 条目 | 弹出详情对话框 |

### 2.7 Entry 详情对话框
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 59 | 在任意页面点击条目打开详情 | 显示完整信息：类型/标题/内容/项目/标签/时间 |
| 60 | 点击 "Copy" 按钮 | 内容复制到剪贴板（粘贴验证） |
| 61 | 点击 "Favorite" 按钮 | 标记为收藏（再次点击取消） |
| 62 | 点击 "Archive" 按钮 | 条目归档（从 Actions/Knowledge 列表消失） |
| 63 | 点击 "Edit" 按钮 | 弹出编辑表单 |
| 64 | 修改 Title 和 Tags，点击确认 | 修改成功，详情刷新 |
| 65 | 点击 "Delete" 按钮 | 条目被删除 |

### 2.8 Dangerous 命令确认
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 66 | 添加一条命令：`rm -rf /tmp/test` | 条目被自动标记为 dangerous |
| 67 | 添加一条命令：`Format-Volume -DriveLetter D` | 条目被自动标记为 dangerous |
| 68 | 在详情中查看 dangerous 条目 | 有危险标记提示 |

### 2.9 Settings 页
| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 69 | 切换到 "Settings" 页 | 显示版本、数据目录、数据库名、条目数 |
| 70 | 检查数据路径 | 显示 `%APPDATA%\DevMemory` 或 portable 路径 |
| 71 | 点击 "Export Today (Markdown)" | 弹出文件保存对话框 |
| 72 | 选择保存路径，确认 | 文件保存成功 |
| 73 | 点击 "Export JSON" | 弹出文件保存对话框 |
| 74 | 保存后检查 JSON 文件 | 内容完整 |
| 75 | 点击 "Import JSON" | 弹出文件选择对话框 |
| 76 | 选择之前导出的 JSON 文件 | 导入成功，条目数增加 |

---

## 三、数据持久化

| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 77 | 添加若干条目后关闭窗口 | 程序正常退出 |
| 78 | 重新启动 `devmemory.exe` | 打开窗口 |
| 79 | 检查 Today/Actions/Knowledge 页 | 之前添加的条目都还在 |
| 80 | CLI 运行 `devmemory.exe list` | 条目数与 GUI 显示一致 |

---

## 四、Portable 模式

| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 81 | 在 exe 同级目录创建 `devmemory-data` 文件夹 | — |
| 82 | 运行 `devmemory.exe` | 使用 `devmemory-data/devmemory.db` 作为数据库 |
| 83 | 添加一条条目 | — |
| 84 | 检查 `devmemory-data/devmemory.db` 存在 | 文件存在且大小 > 0 |
| 85 | 删除 `devmemory-data` 后重新运行 | 回到默认 `%APPDATA%` 路径 |

---

## 五、边界情况

| 步骤 | 操作 | 预期结果 |
|------|------|----------|
| 86 | Content 为空时点击 Save | 给出错误提示，不创建空条目 |
| 87 | 搜索框为空时搜索 | 不崩溃（返回空或全部结果） |
| 88 | 编辑时不修改任何字段直接确认 | 原数据不变 |
| 89 | 导入格式错误的 JSON 文件 | 给出错误提示，不崩溃 |
| 90 | 连续快速点击多个条目 | 不崩溃，最后一次点击的条目显示 |

---

## 测试结果记录

| 测试项 | 通过/失败 | 备注 |
|--------|----------|------|
| 1-2 (启动+版本) | | |
| 3-8 (添加) | | |
| 9-12 (列表) | | |
| 13-16 (搜索) | | |
| 17-18 (查看) | | |
| 19-21 (编辑) | | |
| 22-23 (删除) | | |
| 24 (today) | | |
| 25-30 (导出导入) | | |
| 31-34 (GUI 启动) | | |
| 35-43 (Capture) | | |
| 44-49 (Search) | | |
| 50-52 (Today) | | |
| 53-55 (Actions) | | |
| 56-58 (Knowledge) | | |
| 59-65 (Entry 对话框) | | |
| 66-68 (Dangerous) | | |
| 69-76 (Settings) | | |
| 77-80 (持久化) | | |
| 81-85 (Portable) | | |
| 86-90 (边界) | | |
