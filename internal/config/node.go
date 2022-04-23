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
	NodeJsFlavour = node.NodeJs()
	IoJsFlavour   = node.IoJs()

	NodeHomePathDefault     = "%s" + NodeJsFlavour.Value() + PathSeparator
	NodeCachePathDefault    = "%s" + NodeJsFlavour.Value() + PathSeparator + "cache" + PathSeparator
	NodeVersionsPathDefault = "%s" + NodeJsFlavour.Value() + PathSeparator + "versions" + PathSeparator
	NodeCurrentPathDefault  = "%s" + NodeJsFlavour.Value() + PathSeparator + "current"
)

var NodeDefaultConfig = NodeConfig{
	HomeDir:     NodeHomePathDefault,
	CacheDir:    NodeCachePathDefault,
	VersionsDir: NodeVersionsPathDefault,
	CurrentDir:  NodeCurrentPathDefault,
}
