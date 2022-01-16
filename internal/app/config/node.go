package config

import (
	"fmt"
	"os"

	"github.com/romycode/mvm/internal/node"
	"github.com/romycode/mvm/pkg/env"
)

type Node struct {
	HomeDir     string `json:"home_dir"`
	CacheDir    string `json:"cache_dir"`
	VersionsDir string `json:"versions_dir"`
	CurrentDir  string `json:"current_dir"`
}

var (
	PathSeparator = string(os.PathSeparator)

	NodeJs = node.NodeJs()
	IoJs   = node.IoJs()

	nodeHomePathDefault     = "%s" + NodeJs.Value() + PathSeparator
	nodeCachePathDefault    = "%s" + NodeJs.Value() + PathSeparator + "cache" + PathSeparator
	nodeVersionsPathDefault = "%s" + NodeJs.Value() + PathSeparator + "versions" + PathSeparator
	nodeCurrentPathDefault  = "%s" + NodeJs.Value() + PathSeparator + "current"
)

func populateNodeEnv(mvmHome string) (Node, error) {
	config := Node{}

	config.HomeDir = env.Get("MVM_NODE_HOME", fmt.Sprintf(nodeHomePathDefault, mvmHome))
	if err := os.MkdirAll(config.HomeDir, 0755); err != nil {
		return Node{}, err
	}

	config.CacheDir = env.Get("MVM_NODE_CACHE", fmt.Sprintf(nodeCachePathDefault, mvmHome))
	if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
		return Node{}, err
	}

	config.VersionsDir = env.Get("MVM_NODE_VERSIONS", fmt.Sprintf(nodeVersionsPathDefault, mvmHome))
	if err := os.MkdirAll(config.VersionsDir, 0755); err != nil {
		return Node{}, err
	}

	config.CurrentDir = env.Get("MVM_NODE_CURRENT", fmt.Sprintf(nodeCurrentPathDefault, mvmHome))
	if err := os.MkdirAll(config.CurrentDir, 0755); err != nil {
		return Node{}, err
	}

	return config, nil
}
