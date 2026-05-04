//go:build windows

package config

import (
	"os"
	"path/filepath"
)

func defaultDataDir() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return filepath.Join(".", "devmemory-data")
	}
	return filepath.Join(appData, "DevMemory")
}
