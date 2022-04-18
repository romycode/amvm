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
	DefaultNodeJsFlavour = node.NodeJs()
	IoJsFlavour          = node.IoJs()

	NodeHomePathDefault     = "%s" + DefaultNodeJsFlavour.Value() + PathSeparator
	NodeCachePathDefault    = "%s" + DefaultNodeJsFlavour.Value() + PathSeparator + "cache" + PathSeparator
	NodeVersionsPathDefault = "%s" + DefaultNodeJsFlavour.Value() + PathSeparator + "versions" + PathSeparator
	NodeCurrentPathDefault  = "%s" + DefaultNodeJsFlavour.Value() + PathSeparator + "current"
)

var NodeDefaultConfig = NodeConfig{
	HomeDir:     NodeHomePathDefault,
	CacheDir:    NodeCachePathDefault,
	VersionsDir: NodeVersionsPathDefault,
	CurrentDir:  NodeCurrentPathDefault,
}
