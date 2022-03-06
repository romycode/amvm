package config

import (
	"os"
	"path/filepath"
)

var PathSeparator = string(os.PathSeparator)

var AmvmHomeDirDefault = filepath.Join(os.Getenv("HOME"), ".amvm") + string(os.PathSeparator)

type AmvmConfig struct {
	HomeDir string     `json:"-"`
	Node    NodeConfig `json:"node"`
}
