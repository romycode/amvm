package config

import (
	"github.com/romycode/amvm/internal/pnpm"
)

type PnpmConfig struct {
	HomeDir     string `json:"home_dir"`
	CacheDir    string `json:"cache_dir"`
	VersionsDir string `json:"versions_dir"`
	CurrentDir  string `json:"current_dir"`
}

var (
	PnpmFlavour = pnpm.PnpmJs()

	PnpmHomePathDefault     = "%s" + PnpmFlavour.Value() + PathSeparator
	PnpmCachePathDefault    = "%s" + PnpmFlavour.Value() + PathSeparator + "cache" + PathSeparator
	PnpmVersionsPathDefault = "%s" + PnpmFlavour.Value() + PathSeparator + "versions" + PathSeparator
	PnpmCurrentPathDefault  = "%s" + PnpmFlavour.Value() + PathSeparator + "current"
)

var PnpmDefaultConfig = PnpmConfig{
	HomeDir:     PnpmHomePathDefault,
	CacheDir:    PnpmCachePathDefault,
	VersionsDir: PnpmVersionsPathDefault,
	CurrentDir:  PnpmCurrentPathDefault,
}
