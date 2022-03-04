package config

import (
	"github.com/romycode/amvm/internal/node"
)

type NodeConfig struct {
	HomeDir     string `json:"home_dir"`
	CacheDir    string `json:"cache_dir"`
	VersionsDir string `json:"versions_dir"`
	CurrentDir  string `json:"current_dir"`
}

var (
	DefaultFlavour = node.NodeJs()
	IoJsFlavour    = node.IoJs()

	HomePathDefault     = "%s" + DefaultFlavour.Value() + PathSeparator
	CachePathDefault    = "%s" + DefaultFlavour.Value() + PathSeparator + "cache" + PathSeparator
	VersionsPathDefault = "%s" + DefaultFlavour.Value() + PathSeparator + "versions" + PathSeparator
	CurrentPathDefault  = "%s" + DefaultFlavour.Value() + PathSeparator + "current"
)

var DefaultConfig = NodeConfig{
	HomeDir:     HomePathDefault,
	CacheDir:    CachePathDefault,
	VersionsDir: VersionsPathDefault,
	CurrentDir:  CurrentPathDefault,
}
