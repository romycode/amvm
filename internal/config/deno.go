package config

import (
	"github.com/romycode/amvm/internal/deno"
)

type DenoConfig struct {
	HomeDir     string `json:"home_dir"`
	CacheDir    string `json:"cache_dir"`
	VersionsDir string `json:"versions_dir"`
	CurrentDir  string `json:"current_dir"`
}

var (
	DenoJsFlavour = deno.DenoJs()

	DenoHomePathDefault     = "%sdeno" + PathSeparator
	DenoCachePathDefault    = "%sdeno" + PathSeparator + "cache" + PathSeparator
	DenoVersionsPathDefault = "%sdeno" + PathSeparator + "versions" + PathSeparator
	DenoCurrentPathDefault  = "%sdeno" + PathSeparator + "current"
)

var DenoDefaultConfig = DenoConfig{
	HomeDir:     DenoHomePathDefault,
	CacheDir:    DenoCachePathDefault,
	VersionsDir: DenoVersionsPathDefault,
	CurrentDir:  DenoCurrentPathDefault,
}
