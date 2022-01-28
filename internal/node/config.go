package node

import (
	"os"
)

type Config struct {
	HomeDir     string `json:"home_dir"`
	CacheDir    string `json:"cache_dir"`
	VersionsDir string `json:"versions_dir"`
	CurrentDir  string `json:"current_dir"`
}

var (
	PathSeparator = string(os.PathSeparator)

	DefaultFlavour = NodeJs()
	IoJsFlavour    = IoJs()

	HomePathDefault     = "%s" + DefaultFlavour.Value() + PathSeparator
	CachePathDefault    = "%s" + DefaultFlavour.Value() + PathSeparator + "cache" + PathSeparator
	VersionsPathDefault = "%s" + DefaultFlavour.Value() + PathSeparator + "versions" + PathSeparator
	CurrentPathDefault  = "%s" + DefaultFlavour.Value() + PathSeparator + "current"
)

var DefaultConfig = Config{
	HomeDir:     HomePathDefault,
	CacheDir:    CachePathDefault,
	VersionsDir: VersionsPathDefault,
	CurrentDir:  CurrentPathDefault,
}
