package config

import (
	"github.com/romycode/amvm/internal/java"
)

type JavaConfig struct {
	HomeDir     string `json:"home_dir"`
	CacheDir    string `json:"cache_dir"`
	VersionsDir string `json:"versions_dir"`
	CurrentDir  string `json:"current_dir"`
}

var (
	JavaFlavour = java.Java()

	JavaHomePathDefault     = "%s" + JavaFlavour.Value() + PathSeparator
	JavaCachePathDefault    = "%s" + JavaFlavour.Value() + PathSeparator + "cache" + PathSeparator
	JavaVersionsPathDefault = "%s" + JavaFlavour.Value() + PathSeparator + "versions" + PathSeparator
	JavaCurrentPathDefault  = "%s" + JavaFlavour.Value() + PathSeparator + "current"
)

var JavaDefaultConfig = JavaConfig{
	HomeDir:     JavaHomePathDefault,
	CacheDir:    JavaCachePathDefault,
	VersionsDir: JavaVersionsPathDefault,
	CurrentDir:  JavaCurrentPathDefault,
}
