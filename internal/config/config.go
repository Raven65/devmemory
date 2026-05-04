package config

// Config holds application configuration.
type Config struct {
	DataDir string `json:"data_dir"`
	Port    int    `json:"port"`
	DBPath  string `json:"-"`
}

const (
	DefaultPort = 8420
	DBFileName  = "devmemory.db"
	AppName     = "devmemory"
)
