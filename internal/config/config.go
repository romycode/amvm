package config

import (
	"os"
	"path/filepath"
)

var PathSeparator = string(os.PathSeparator)

var MvmHomeDirDefault = filepath.Join(os.Getenv("HOME"), ".mvm") + string(os.PathSeparator)

type MvmConfig struct {
	HomeDir string     `json:"-"`
	Node    NodeConfig `json:"node"`
}
