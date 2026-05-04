# DevMemory

A local-first personal memory hub, action library, and work journal for developers.

## What It Does

DevMemory helps you capture, search, and reuse the碎片信息 that accumulates during daily development:

- Shell commands
- Useful URLs
- Code snippets & AI prompts
- Problem-solving notes
- Business knowledge
- Daily work logs

Core loop: **Capture → Classify → Search → Act → Export**

## Quick Start

### Build

```bash
go build -o devmemory ./cmd/devmemory
```

### Usage

```bash
# Add entries
devmemory add "ss -tinp | grep ESTAB" --type command --tags linux,tcp,debug
devmemory add "https://docs.kernel.org/networking/" --type url --tags linux,kernel
devmemory note "NUMA binding affects RDMA performance" --project rdma --tags numa,issue

# List entries
devmemory list
devmemory list --type command
devmemory list --tag linux

# Search
devmemory search "tcp estab"
devmemory search "numa rdma"

# Today's entries
devmemory today

# Version
devmemory version
```

### Auto Type Detection

If `--type` is not specified, DevMemory infers the type:

- Starts with `http://` or `https://` → `url`
- Contains `|`, `&&`, `||`, `;`, `$` → `command`
- Otherwise → `note`

## Cross-Platform Build

```bash
# From Linux, build all platforms:
./scripts/build.sh

# Or manually:
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/devmemory-linux-amd64 ./cmd/devmemory
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/devmemory-windows-amd64.exe ./cmd/devmemory
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/devmemory-darwin-amd64 ./cmd/devmemory
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/devmemory-darwin-arm64 ./cmd/devmemory
```

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

- **Language**: Go
- **Storage**: bbolt (embedded key-value)
- **Server**: net/http (stdlib)
- **Frontend**: HTML/CSS/JS (Go embed)
- **CLI**: stdlib flag + manual subcommand dispatch

## Design Principles

- Local-first, no cloud dependency
- Single binary, zero install
- Capture first, organize later
- Personal use, no multi-user concerns
- Minimal dependencies, no CGO
