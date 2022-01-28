package config

import (
	"github.com/romycode/mvm/internal/node"
	"os"
	"path/filepath"
)

type MvmConfig struct {
	HomeDir string      `json:"-"`
	Node    node.Config `json:"node"`
}

var MvmHomeDirDefault = filepath.Join(os.Getenv("HOME"), ".mvm") + string(os.PathSeparator)
