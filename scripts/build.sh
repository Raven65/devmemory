#!/usr/bin/env bash
set -euo pipefail

APP_NAME="devmemory"
VERSION="${1:-dev}"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
DIST_DIR="dist"

LDFLAGS="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${DATE}"

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

echo "Building ${APP_NAME} ${VERSION}..."

echo "  Linux amd64..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-linux-amd64" ./cmd/devmemory

echo "  Windows amd64..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-windows-amd64.exe" ./cmd/devmemory

echo "  macOS amd64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-darwin-amd64" ./cmd/devmemory

echo "  macOS arm64..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-darwin-arm64" ./cmd/devmemory

echo ""
echo "Build complete:"
ls -lh "${DIST_DIR}/"
