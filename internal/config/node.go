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
	NodeFlavour = node.NodeJs()

	NodeHomePathDefault     = "%s" + NodeFlavour.Value() + PathSeparator
	NodeCachePathDefault    = "%s" + NodeFlavour.Value() + PathSeparator + "cache" + PathSeparator
	NodeVersionsPathDefault = "%s" + NodeFlavour.Value() + PathSeparator + "versions" + PathSeparator
	NodeCurrentPathDefault  = "%s" + NodeFlavour.Value() + PathSeparator + "current"
)

var NodeDefaultConfig = NodeConfig{
	HomeDir:     NodeHomePathDefault,
	CacheDir:    NodeCachePathDefault,
	VersionsDir: NodeVersionsPathDefault,
	CurrentDir:  NodeCurrentPathDefault,
}
