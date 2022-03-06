package main

import (
	"encoding/json"
	"fmt"
	httpstd "net/http"
	"os"
	"path/filepath"

	"github.com/romycode/amvm/internal/app/cmd"
	"github.com/romycode/amvm/internal/app/fetch"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/pkg/color"
	"github.com/romycode/amvm/pkg/env"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/http"
)

type Command string

const (
	Info    Command = "info"
	Fetch   Command = "fetch"
	Install Command = "install"
	Use     Command = "use"
)

func createDefaultConfigIfIsNecessary(path string) error {
	if !file.Exists(path) {
		data, _ := json.Marshal(config.MvmConfig{
			HomeDir: config.MvmHomeDirDefault,
			Node:    config.DefaultConfig,
		})

		if err := file.Write(path, data); err != nil {
			return fmt.Errorf("error creating default configuration file: %s", path)
		}
	}
	return nil
}

func readConfig(path string) (*config.MvmConfig, error) {
	data, err := file.Read(path)
	if err != nil {
		return &config.MvmConfig{}, fmt.Errorf("error reading configuration file: %s", path)
	}

	c := &config.MvmConfig{HomeDir: filepath.Dir(path)}
	err = json.Unmarshal(data, c)

	if err != nil {
		return &config.MvmConfig{}, fmt.Errorf("invalid configuration file: %s", path)
	}

	return c, nil
}

func writeConfig(path string, config config.MvmConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err = file.Write(path, data); err != nil {
		return err
	}

	return nil
}

func loadNodeConfig(mvmHome string) (config.NodeConfig, error) {
	c := config.NodeConfig{}

	c.HomeDir = env.Get("MVM_NODE_HOME", fmt.Sprintf(config.HomePathDefault, mvmHome))
	if err := os.MkdirAll(c.HomeDir, 0755); err != nil && !file.Exists(c.HomeDir) {
		return config.NodeConfig{}, err
	}

	c.CacheDir = env.Get("MVM_NODE_CACHE", fmt.Sprintf(config.CachePathDefault, mvmHome))
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil && !file.Exists(c.CacheDir) {
		return config.NodeConfig{}, err
	}

	c.VersionsDir = env.Get("MVM_NODE_VERSIONS", fmt.Sprintf(config.VersionsPathDefault, mvmHome))
	if err := os.MkdirAll(c.VersionsDir, 0755); err != nil && !file.Exists(c.VersionsDir) {
		return config.NodeConfig{}, err
	}

	c.CurrentDir = env.Get("MVM_NODE_CURRENT", fmt.Sprintf(config.CurrentPathDefault, mvmHome))

	return c, nil
}

func loadConfiguration() (*config.MvmConfig, error) {
	mvmPath := env.Get("MVM_HOME", config.MvmHomeDirDefault)
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
	if err := writeConfig(configFilePath, *c); err != nil {
		return nil, err
	}

	return c, nil
}

func PrintOutput(output cmd.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}

func main() {
	conf, err := loadConfiguration()
	if err != nil {
		PrintOutput(cmd.NewOutput(color.Colorize(err.Error(), color.Red), 1))
	}
	if 1 == len(os.Args) {
		PrintOutput(cmd.NewOutput(color.Colorize("use: amvm <info|install|use|fetch> <nodejs> <flavour> <version>", color.Green), 0))
	}

	nhc := http.NewClient(httpstd.DefaultClient, fetch.NodeJsURLTemplate)
	nf := fetch.NewNodeJsFetcher(nhc)

	command := Command(os.Args[1])
	switch command {
	case Info:
		PrintOutput(cmd.NewInfoCommand(nf).Run())
	case Fetch:
		PrintOutput(cmd.NewFetchCommand(conf, nf).Run())
	case Install:
		PrintOutput(cmd.NewInstallCommand(conf, nf, nhc).Run())
	case Use:
		PrintOutput(cmd.NewUseCommand(conf, nf).Run())
	}
}
