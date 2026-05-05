#!/usr/bin/env bash
set -euo pipefail

APP_NAME="devmemory"
VERSION="${1:-dev}"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
DIST_DIR="dist"

LDFLAGS="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${DATE}"
WIN_LDFLAGS="-s -w -H windowsgui -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Date=${DATE}"

# --- Platform detection helpers ---
has_cmd() { command -v "$1" &>/dev/null; }

# mingw-w64 for Windows cross-compilation
has_mingw() { command -v x86_64-w64-mingw32-gcc &>/dev/null; }

# fyne-cross (Docker-based, supports all platforms)
has_fyne_cross() { command -v fyne-cross &>/dev/null && has_cmd docker; }

# zig as drop-in C cross-compiler
has_zig() { command -v zig &>/dev/null; }

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

echo "=== DevMemory Build Script ==="
echo "Version: ${VERSION}  Commit: ${COMMIT}"
echo ""

# ---- Native build ----
build_native() {
	echo "[1/4] Native build (current platform)..."
	CGO_ENABLED=1 go build -trimpath -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}" ./cmd/devmemory
	echo "  OK: ${DIST_DIR}/${APP_NAME} ($(du -h "${DIST_DIR}/${APP_NAME}" | cut -f1))"
}

# ---- Windows amd64 ----
build_windows() {
	echo "[2/4] Windows amd64..."

	if has_mingw; then
		echo "  Using mingw-w64..."
		CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
			go build -trimpath -ldflags "${WIN_LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-windows-amd64.exe" ./cmd/devmemory
		echo "  OK: ${DIST_DIR}/${APP_NAME}-windows-amd64.exe ($(du -h "${DIST_DIR}/${APP_NAME}-windows-amd64.exe" | cut -f1))"
	elif has_zig; then
		echo "  Using zig cc..."
		CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" \
			go build -trimpath -ldflags "${WIN_LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-windows-amd64.exe" ./cmd/devmemory
		echo "  OK: ${DIST_DIR}/${APP_NAME}-windows-amd64.exe ($(du -h "${DIST_DIR}/${APP_NAME}-windows-amd64.exe" | cut -f1))"
	else
		echo "  SKIP: no mingw-w64 or zig found. Install one of:"
		echo "    • dnf install mingw64-gcc"
		echo "    • apt install gcc-mingw-w64-x86-64"
		echo "    • Download zig from https://ziglang.org/download/"
	fi
}

# ---- Darwin amd64 ----
build_darwin_amd64() {
	echo "[3/4] Darwin amd64..."

	if has_zig; then
		echo "  Using zig cc..."
		CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC="zig cc -target x86_64-macos" \
			go build -trimpath -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-darwin-amd64" ./cmd/devmemory
		echo "  OK: ${DIST_DIR}/${APP_NAME}-darwin-amd64 ($(du -h "${DIST_DIR}/${APP_NAME}-darwin-amd64" | cut -f1))"
	else
		echo "  SKIP: no zig found. Options:"
		echo "    • Build on macOS directly"
		echo "    • Use fyne-cross with Docker"
		echo "    • Use GitHub Actions (macOS runner)"
	fi
}

# ---- Darwin arm64 ----
build_darwin_arm64() {
	echo "[4/4] Darwin arm64..."

	if has_zig; then
		echo "  Using zig cc..."
		CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="zig cc -target aarch64-macos" \
			go build -trimpath -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${APP_NAME}-darwin-arm64" ./cmd/devmemory
		echo "  OK: ${DIST_DIR}/${APP_NAME}-darwin-arm64 ($(du -h "${DIST_DIR}/${APP_NAME}-darwin-arm64" | cut -f1))"
	else
		echo "  SKIP: no zig found. Options:"
		echo "    • Build on Apple Silicon Mac directly"
		echo "    • Use fyne-cross with Docker"
		echo "    • Use GitHub Actions (macOS runner)"
	fi
}

# ---- Docker-based build (all platforms at once) ----
build_fyne_cross() {
	echo "=== Using fyne-cross (Docker) ==="
	fyne-cross --targets=linux/amd64,windows/amd64,darwin/amd64,darwin/arm64 \
		-ldflags "${LDFLAGS}" -app-version "${VERSION}" .
	echo ""
	echo "Outputs in fyne-cross/dist/"
}

# ---- Main dispatch ----
if has_fyne_cross && [[ "${USE_FYNE_CROSS:-}" == "1" ]]; then
	build_fyne_cross
else
	build_native
	echo ""
	build_windows
	echo ""
	build_darwin_amd64
	echo ""
	build_darwin_arm64
fi

echo ""
echo "=== Build artifacts ==="
ls -lh "${DIST_DIR}/" 2>/dev/null || true
echo ""
echo "=== Done ==="
