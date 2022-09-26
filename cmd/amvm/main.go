package main

import (
	"encoding/json"
	"fmt"
	httpstd "net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/romycode/amvm/internal/app/fetch"
	"github.com/romycode/amvm/pkg/http"
	"github.com/romycode/amvm/pkg/ui"

	"github.com/romycode/amvm/internal/app/cmd"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/pkg/env"
	"github.com/romycode/amvm/pkg/file"
)

type Command string

const (
	Info    Command = "info"
	Fetch   Command = "fetch"
	Install Command = "install"
	Use     Command = "use"
)

func main() {
	conf, err := loadConfiguration()
	if err != nil {
		PrintOutput(cmd.NewOutput(err.Error(), ui.Red, 1))
	}
	if 1 == len(os.Args) {
		PrintOutput(cmd.NewOutput("use: amvm <info|install|use|fetch> <nodejs> <flavour> <version>", ui.Green, 0))
	}

	hc := http.NewClient(&httpstd.Client{}, "")
	nhc := http.NewClient(&httpstd.Client{}, fetch.NodeJsURLTemplate)
	dhc := http.NewClient(&httpstd.Client{}, fetch.DenoGithubURLTemplate)
	phc := http.NewClient(&httpstd.Client{}, fetch.PnpmJsURLTemplate)
	jhc := http.NewClient(&httpstd.Client{}, fetch.JavaURLApi)

	arch := runtime.GOARCH
	system := runtime.GOOS

	nf := fetch.NewNodeJsFetcher(nhc, arch, system)
	df := fetch.NewDenoFetcher(dhc, arch, system)
	pf := fetch.NewPnpmJsFetcher(phc, arch, system)
	jf := fetch.NewJavaFetcher(jhc, arch, system)
	ff := fetch.NewFactory(nf, df, pf, jf)

	command := Command(os.Args[1])
	switch command {
	case Info:
		PrintOutput(cmd.NewInfoCommand(ff).Run())
	case Fetch:
		PrintOutput(cmd.NewFetchCommand(conf, ff).Run())
	case Install:
		PrintOutput(cmd.NewInstallCommand(conf, ff, hc).Run())
	case Use:
		PrintOutput(cmd.NewUseCommand(conf, ff).Run())
	}
}

func createDefaultConfigIfIsNecessary(path string) error {
	if !file.Exists(path) {
		data, _ := json.Marshal(
			config.AmvmConfig{
				HomeDir: config.AmvmHomeDirDefault,
				Node:    config.NodeDefaultConfig,
				Deno:    config.DenoDefaultConfig,
				Pnpm:    config.PnpmDefaultConfig,
				Java:    config.JavaDefaultConfig,
			},
		)

		if err := file.Write(path, data); err != nil {
			return fmt.Errorf("error creating default configuration file: %s", path)
		}
	}
	return nil
}

func loadConfiguration() (*config.AmvmConfig, error) {
	mvmPath := env.Get("AMVM_HOME", config.AmvmHomeDirDefault)
	if err := os.MkdirAll(mvmPath, 0755); err != nil && !file.Exists(mvmPath) {
		return nil, err
	}

	configFilePath := fmt.Sprintf("%sconfig.json", mvmPath)
	if err := createDefaultConfigIfIsNecessary(configFilePath); err != nil {
		return nil, err
	}

	c, err := readConfig(configFilePath)
	if err != nil {
		return nil, err
	}

	if c.Node, err = loadNodeConfig(mvmPath); err != nil {
		return nil, err
	}
	if c.Deno, err = loadDenoConfig(mvmPath); err != nil {
		return nil, err
	}
	if c.Pnpm, err = loadPnpmConfig(mvmPath); err != nil {
		return nil, err
	}
	if c.Java, err = loadJavaConfig(mvmPath); err != nil {
		return nil, err
	}
	if err := writeConfig(configFilePath, *c); err != nil {
		return nil, err
	}

	return c, nil
}

func loadNodeConfig(mvmHome string) (config.NodeConfig, error) {
	c := config.NodeConfig{}

	c.HomeDir = env.Get("AMVM_NODE_HOME", fmt.Sprintf(config.NodeHomePathDefault, mvmHome))
	if err := os.MkdirAll(c.HomeDir, 0755); err != nil && !file.Exists(c.HomeDir) {
		return config.NodeConfig{}, err
	}

	c.CacheDir = env.Get("AMVM_NODE_CACHE", fmt.Sprintf(config.NodeCachePathDefault, mvmHome))
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil && !file.Exists(c.CacheDir) {
		return config.NodeConfig{}, err
	}

	c.VersionsDir = env.Get("AMVM_NODE_VERSIONS", fmt.Sprintf(config.NodeVersionsPathDefault, mvmHome))
	if err := os.MkdirAll(c.VersionsDir, 0755); err != nil && !file.Exists(c.VersionsDir) {
		return config.NodeConfig{}, err
	}

	c.CurrentDir = env.Get("AMVM_NODE_CURRENT", fmt.Sprintf(config.NodeCurrentPathDefault, mvmHome))

	return c, nil
}

func loadDenoConfig(mvmHome string) (config.DenoConfig, error) {
	c := config.DenoConfig{}

	c.HomeDir = env.Get("AMVM_DENO_HOME", fmt.Sprintf(config.DenoHomePathDefault, mvmHome))
	if err := os.MkdirAll(c.HomeDir, 0755); err != nil && !file.Exists(c.HomeDir) {
		return config.DenoConfig{}, err
	}

	c.CacheDir = env.Get("AMVM_DENO_CACHE", fmt.Sprintf(config.DenoCachePathDefault, mvmHome))
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil && !file.Exists(c.CacheDir) {
		return config.DenoConfig{}, err
	}

	c.VersionsDir = env.Get("AMVM_DENO_VERSIONS", fmt.Sprintf(config.DenoVersionsPathDefault, mvmHome))
	if err := os.MkdirAll(c.VersionsDir, 0755); err != nil && !file.Exists(c.VersionsDir) {
		return config.DenoConfig{}, err
	}

	c.CurrentDir = env.Get("AMVM_DENO_CURRENT", fmt.Sprintf(config.DenoCurrentPathDefault, mvmHome))

	return c, nil
}

func loadPnpmConfig(mvmHome string) (config.PnpmConfig, error) {
	c := config.PnpmConfig{}

	c.HomeDir = env.Get("AMVM_PNPM_HOME", fmt.Sprintf(config.PnpmHomePathDefault, mvmHome))
	if err := os.MkdirAll(c.HomeDir, 0755); err != nil && !file.Exists(c.HomeDir) {
		return config.PnpmConfig{}, err
	}

	c.CacheDir = env.Get("AMVM_PNPM_CACHE", fmt.Sprintf(config.PnpmCachePathDefault, mvmHome))
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil && !file.Exists(c.CacheDir) {
		return config.PnpmConfig{}, err
	}

	c.VersionsDir = env.Get("AMVM_PNPM_VERSIONS", fmt.Sprintf(config.PnpmVersionsPathDefault, mvmHome))
	if err := os.MkdirAll(c.VersionsDir, 0755); err != nil && !file.Exists(c.VersionsDir) {
		return config.PnpmConfig{}, err
	}

	c.CurrentDir = env.Get("AMVM_PNPM_CURRENT", fmt.Sprintf(config.PnpmCurrentPathDefault, mvmHome))

	return c, nil
}

func loadJavaConfig(mvmHome string) (config.JavaConfig, error) {
	c := config.JavaConfig{}

	c.HomeDir = env.Get("AMVM_JAVA_HOME", fmt.Sprintf(config.JavaHomePathDefault, mvmHome))
	if err := os.MkdirAll(c.HomeDir, 0755); err != nil && !file.Exists(c.HomeDir) {
		return config.JavaConfig{}, err
	}

	c.CacheDir = env.Get("AMVM_JAVA_CACHE", fmt.Sprintf(config.JavaCachePathDefault, mvmHome))
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil && !file.Exists(c.CacheDir) {
		return config.JavaConfig{}, err
	}

	c.VersionsDir = env.Get("AMVM_JAVA_VERSIONS", fmt.Sprintf(config.JavaVersionsPathDefault, mvmHome))
	if err := os.MkdirAll(c.VersionsDir, 0755); err != nil && !file.Exists(c.VersionsDir) {
		return config.JavaConfig{}, err
	}

	c.CurrentDir = env.Get("AMVM_JAVA_CURRENT", fmt.Sprintf(config.JavaCurrentPathDefault, mvmHome))

	return c, nil
}

func readConfig(path string) (*config.AmvmConfig, error) {
	data, err := file.Read(path)
	if err != nil {
		return &config.AmvmConfig{}, fmt.Errorf("error reading configuration file: %s", path)
	}

	c := &config.AmvmConfig{HomeDir: filepath.Dir(path)}
	err = json.Unmarshal(data, c)

	if err != nil {
		return &config.AmvmConfig{}, fmt.Errorf("invalid configuration file: %s", path)
	}

	return c, nil
}

func writeConfig(path string, config config.AmvmConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err = file.Write(path, data); err != nil {
		return err
	}

	return nil
}

func PrintOutput(output cmd.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}
