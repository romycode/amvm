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
	DefaultPnpmJsFlavour = pnpm.PnpmJs()

	PnpmHomePathDefault     = "%s" + DefaultPnpmJsFlavour.Value() + PathSeparator
	PnpmCachePathDefault    = "%s" + DefaultPnpmJsFlavour.Value() + PathSeparator + "cache" + PathSeparator
	PnpmVersionsPathDefault = "%s" + DefaultPnpmJsFlavour.Value() + PathSeparator + "versions" + PathSeparator
	PnpmCurrentPathDefault  = "%s" + DefaultPnpmJsFlavour.Value() + PathSeparator + "current"
)

var PnpmDefaultConfig = PnpmConfig{
	HomeDir:     PnpmHomePathDefault,
	CacheDir:    PnpmCachePathDefault,
	VersionsDir: PnpmVersionsPathDefault,
	CurrentDir:  PnpmCurrentPathDefault,
}
