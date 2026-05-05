# DevMemory

A local-first personal memory hub, action library, and work journal for developers.

## What It Does

DevMemory helps you capture, search, and reuse the bits of information that accumulates during daily development:

- Shell commands
- Useful URLs
- Code snippets & AI prompts
- Problem-solving notes
- Business knowledge
- Daily work logs

Core loop: **Capture → Classify → Search → Act → Export**

## Quick Start

### Build

Requirements: Go 1.23+, GCC (CGO for Fyne GUI)

```bash
go build -o devmemory ./cmd/devmemory
```

### Usage

**Desktop GUI** (default, no arguments):

```bash
./devmemory          # Launch Fyne desktop GUI
./devmemory gui      # Same as above
```

**CLI commands:**

```bash
# Add entries
devmemory add "ss -tinp | grep ESTAB" --type command --tags linux,tcp,debug
devmemory add "https://docs.kernel.org/networking/" --type url --tags linux,kernel
devmemory add "NUMA binding affects RDMA performance" --type note --project rdma --tags numa,issue

# List entries
devmemory list
devmemory list --type command
devmemory list --tag linux

# Search
devmemory search "tcp estab"
devmemory search "numa rdma"

# Today's entries
devmemory today

# Show / Edit / Delete
devmemory show <id>
devmemory edit <id> --title "New Title"
devmemory delete <id>

# Export / Import
devmemory export today -o report.md
devmemory export json -o backup.json
devmemory import backup.json

# Web UI (legacy, still available)
devmemory serve --port 8420

# Version
devmemory version
```

### Auto Type Detection

If `--type` is not specified, DevMemory infers the type:

- Starts with `http://` or `https://` → `url`
- Contains `|`, `&&`, `||`, `;`, `$` → `command`
- Otherwise → `note`

## Desktop GUI

DevMemory uses [Fyne](https://fyne.io/) for a native desktop interface with 6 pages:

| Page | Description |
|------|-------------|
| Capture | Create new entries with type/title/project/tags |
| Search | Full-text search with type filter and relevance scoring |
| Today | Today's entries timeline |
| Actions | Action-type entries (commands, URLs, snippets, prompts, files, folders) |
| Knowledge | Knowledge-type entries (notes, issues, business, journals, tasks) |
| Settings | Data info, export/import, status bar |

## Cross-Platform Build

Fyne requires CGO. For native builds:

```bash
./scripts/build.sh
```

For cross-platform builds, use [fyne-cross](https://github.com/fyne-io/fyne-cross):

```bash
go install github.com/fyne-io/fyne-cross@latest
fyne-cross --targets=linux/amd64,windows/amd64,darwin/amd64,darwin/arm64 .
```

Or build on each target platform directly.

## Data Storage

Default data directory:

| Platform | Path |
|----------|------|
| Linux | `~/.config/devmemory` |
| macOS | `~/Library/Application Support/DevMemory` |
| Windows | `%APPDATA%\DevMemory` |

### Portable Mode

Create a `devmemory-data/` directory next to the binary:

```
devmemory
devmemory-data/
  devmemory.db
```

## Entry Types

| Type | Default Action |
|------|---------------|
| command | Copy or execute (with confirmation) |
| url | Open in browser |
| snippet | Copy text |
| prompt | Copy text |
| note | View |
| issue | View |
| journal | View |
| business | View |
| task | View / mark complete |
| file | Open file |
| folder | Open folder |

## Tech Stack

- **Language**: Go 1.23+
- **GUI**: [Fyne](https://fyne.io/) v2 (native desktop)
- **Storage**: bbolt (embedded key-value)
- **Server**: net/http (stdlib, legacy Web UI)
- **CLI**: stdlib flag + manual subcommand dispatch
- **CGO**: Required (Fyne dependency)

## Design Principles

- Local-first, no cloud dependency
- Single binary, zero install
- Capture first, organize later
- Personal use, no multi-user concerns
- CLI and GUI share the same service layer
