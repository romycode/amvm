package main

import (
	"encoding/json"
	"fmt"
	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
	"github.com/romycode/mvm/internal/node"
	"github.com/romycode/mvm/pkg/env"
	"github.com/romycode/mvm/pkg/file"
	httpstd "net/http"
	"os"
	"path/filepath"

	"github.com/romycode/mvm/internal/app/cmd"
	"github.com/romycode/mvm/pkg/color"
	"github.com/romycode/mvm/pkg/http"
)

type Command string

const (
	Info    Command = "info"
	Fetch   Command = "fetch"
	Install Command = "install"
	Use     Command = "use"
)

func createDefaultConfigIfIsNecessary(path string) error {
	if !file.Check(path) {
		data, _ := json.Marshal(config.MvmConfig{
			HomeDir: config.MvmHomeDirDefault,
			Node:    node.DefaultConfig,
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

func loadNodeConfig(mvmHome string) (node.Config, error) {
	c := node.Config{}

	c.HomeDir = env.Get("MVM_NODE_HOME", fmt.Sprintf(node.HomePathDefault, mvmHome))
	if err := os.MkdirAll(c.HomeDir, 0755); err != nil {
		return node.Config{}, err
	}

	c.CacheDir = env.Get("MVM_NODE_CACHE", fmt.Sprintf(node.CachePathDefault, mvmHome))
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil {
		return node.Config{}, err
	}

	c.VersionsDir = env.Get("MVM_NODE_VERSIONS", fmt.Sprintf(node.VersionsPathDefault, mvmHome))
	if err := os.MkdirAll(c.VersionsDir, 0755); err != nil {
		return node.Config{}, err
	}

	c.CurrentDir = env.Get("MVM_NODE_CURRENT", fmt.Sprintf(node.CurrentPathDefault, mvmHome))
	if err := os.MkdirAll(c.CurrentDir, 0755); err != nil {
		return node.Config{}, err
	}

	return c, nil
}

func loadConfiguration() (*config.MvmConfig, error) {
	mvmPath := env.Get("MVM_HOME", config.MvmHomeDirDefault)
	if err := os.MkdirAll(mvmPath, 0755); err != nil {
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

func main() {
	conf, err := loadConfiguration()
	if err != nil {
		PrintOutput(cmd.NewOutput(color.Colorize(err.Error(), color.Red), 1))
	}
	if 1 == len(os.Args) {
		PrintOutput(cmd.NewOutput(color.Colorize("use: mvm <info|install|use|fetch> <nodejs> <flavour> <version>", color.White), 0))
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
		PrintOutput(cmd.NewInstallCommand(conf, nf).Run())
	case Use:
		PrintOutput(cmd.NewUseCommand(conf, nf).Run())
	}
}

func PrintOutput(output cmd.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}
