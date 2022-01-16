package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/romycode/mvm/pkg/env"
	"github.com/romycode/mvm/pkg/file"
)

type MvmConfig struct {
	HomeDir string `json:"-"`
	Node    Node   `json:"node"`
}

var mvmHomeDirDefault = filepath.Join(os.Getenv("HOME"), ".mvm") + string(os.PathSeparator)

func createDefaultConfigIfIsNecessary(path string) error {
	if !file.Check(path) {
		data, _ := json.Marshal(MvmConfig{
			HomeDir: mvmHomeDirDefault,
			Node: Node{
				HomeDir:     nodeHomePathDefault,
				CacheDir:    nodeCachePathDefault,
				VersionsDir: nodeVersionsPathDefault,
				CurrentDir:  nodeCurrentPathDefault,
			},
		})

		if err := file.Write(path, data); err != nil {
			return fmt.Errorf("error creating default configuration file: %s", path)
		}
	}
	return nil
}

func readConfig(path string) (*MvmConfig, error) {
	data, err := file.Read(path)
	if err != nil {
		return &MvmConfig{}, fmt.Errorf("error reading configuration file: %s", path)
	}

	var config = &MvmConfig{HomeDir: filepath.Dir(path)}
	err = json.Unmarshal(data, config)

	if err != nil {
		return &MvmConfig{}, fmt.Errorf("invalid configuration file: %s", path)
	}

	return config, nil
}

func writeConfig(path string, config MvmConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err = file.Write(path, data); err != nil {
		return err
	}

	return nil
}

func LoadConfiguration() (*MvmConfig, error) {
	mvmPath := env.Get("MVM_HOME", mvmHomeDirDefault)
	if err := os.MkdirAll(mvmPath, 0755); err != nil {
		return nil, err
	}

	configFilePath := fmt.Sprintf("%sconfig.json", mvmPath)

	if err := createDefaultConfigIfIsNecessary(configFilePath); err != nil {
		return nil, err
	}

	config, err := readConfig(configFilePath)
	if err != nil {
		return nil, err
	}

	if config.Node, err = populateNodeEnv(mvmPath); err != nil {
		return nil, err
	}
	if err := writeConfig(configFilePath, *config); err != nil {
		return nil, err
	}

	return config, nil
}
