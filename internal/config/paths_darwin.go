//go:build darwin

package config

import (
	"os"
	"path/filepath"
)

func defaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "devmemory-data")
	}
	return filepath.Join(home, "Library", "Application Support", "DevMemory")
}
