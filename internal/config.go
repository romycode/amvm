package internal

import (
	"os"
	"path/filepath"
)

var (
	PathSeparator      = string(os.PathSeparator)
	AmvmHomeDefaultDir = filepath.Join(os.Getenv("HOME"), ".amvm")
)

type Tool string

const (
	Node Tool = "node"
	Deno Tool = "deno"
	Pnpm Tool = "pnpm"
	Java Tool = "java"
)

type Tools []Tool

var AvailableTools = Tools{Node, Deno, Pnpm, Java}

type Config struct {
	HomeDir     string `json:"home_dir"`
	CacheDir    string `json:"cache_dir"`
	VersionsDir string `json:"versions_dir"`
	CurrentDir  string `json:"current_dir"`
}

type AmvmConfig struct {
	HomeDir string          `json:"-"`
	Tools   map[Tool]Config `json:"tools"`
}

func CreateDefaultConfig(dir string) map[Tool]Config {
	config := map[Tool]Config{}

	for _, tool := range AvailableTools {
		config[tool] = Config{
			HomeDir:     filepath.Join(dir, string(tool)),
			CacheDir:    filepath.Join(dir, string(tool), "cache"),
			CurrentDir:  filepath.Join(dir, string(tool), "current"),
			VersionsDir: filepath.Join(dir, string(tool), "versions"),
		}
	}

	return config
}
