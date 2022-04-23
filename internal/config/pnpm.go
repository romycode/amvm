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
	PnpmJsFlavour = pnpm.PnpmJs()

	PnpmHomePathDefault     = "%s" + PnpmJsFlavour.Value() + PathSeparator
	PnpmCachePathDefault    = "%s" + PnpmJsFlavour.Value() + PathSeparator + "cache" + PathSeparator
	PnpmVersionsPathDefault = "%s" + PnpmJsFlavour.Value() + PathSeparator + "versions" + PathSeparator
	PnpmCurrentPathDefault  = "%s" + PnpmJsFlavour.Value() + PathSeparator + "current"
)

var PnpmDefaultConfig = PnpmConfig{
	HomeDir:     PnpmHomePathDefault,
	CacheDir:    PnpmCachePathDefault,
	VersionsDir: PnpmVersionsPathDefault,
	CurrentDir:  PnpmCurrentPathDefault,
}
