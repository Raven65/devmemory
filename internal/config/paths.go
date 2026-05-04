package config

import (
	"os"
	"path/filepath"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	dir := DataDir()
	return &Config{
		DataDir: dir,
		Port:    DefaultPort,
		DBPath:  filepath.Join(dir, DBFileName),
	}
}

// DataDir returns the data directory path.
// Portable mode (devmemory-data/ next to binary) takes priority.
func DataDir() string {
	// Portable mode: check for devmemory-data/ next to executable
	exePath, err := os.Executable()
	if err == nil {
		portableDir := filepath.Join(filepath.Dir(exePath), "devmemory-data")
		if info, err := os.Stat(portableDir); err == nil && info.IsDir() {
			return portableDir
		}
	}

	return defaultDataDir()
}

// EnsureDataDir creates the data directory if it doesn't exist.
func EnsureDataDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
