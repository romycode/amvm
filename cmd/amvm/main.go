package main

import (
	"encoding/json"
	"fmt"
	httpstd "net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/app/cmd"
	"github.com/romycode/amvm/internal/fetch"
	fstrategies "github.com/romycode/amvm/internal/fetch/strategies"
	"github.com/romycode/amvm/internal/install"
	istrategies "github.com/romycode/amvm/internal/install/strategies"
	"github.com/romycode/amvm/pkg/env"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/http"
	"github.com/romycode/amvm/pkg/ui"
)

type Command string

const (
	Use     Command = "use"
	Info    Command = "info"
	Fetch   Command = "fetch"
	Install Command = "install"
)

var (
	err error
	c   *internal.AmvmConfig
)

func init() {
	c, err = loadConfiguration()
	if err != nil {
		PrintOutput(internal.NewOutput(err.Error(), ui.Red, 1))
	}
}

func main() {
	if len(os.Args) == 1 {
		PrintOutput(internal.NewOutput("use: amvm <info|install|use|fetch> <nodejs> <flavour> <version>", ui.Green, 0))
	}

	arch := runtime.GOARCH
	system := runtime.GOOS

	hc := http.NewClient(&httpstd.Client{})
	nfs := fstrategies.NewNodeJsFetcherStrategy(hc, arch, system)
	pfs := fstrategies.NewPnpmJsFetcherStrategy(hc, arch, system)
	dfs := fstrategies.NewDenoFetcherStrategy(hc, arch, system)
	jfs := fstrategies.NewJavaFetcherStrategy(hc, arch, system)
	f := fetch.NewFetcher([]fetch.Strategy{nfs, pfs, dfs, jfs})

	nis := istrategies.NewNodeJsInstallerStrategy(hc, c.Tools[internal.Node], arch, system)
	pis := istrategies.NewPnpmJsInstallerStrategy(hc, c.Tools[internal.Pnpm], arch, system)
	dis := istrategies.NewDenoInstallerStrategy(hc, c.Tools[internal.Deno], arch, system)
	jis := istrategies.NewJavaInstallerStrategy(hc, c.Tools[internal.Java], arch, system)
	i := install.NewInstaller([]install.Strategy{nis, pis, dis, jis})

	command := Command(os.Args[1])
	switch command {
	case Info:
		PrintOutput(cmd.NewInfoCommand(f).Run())
	case Fetch:
		PrintOutput(cmd.NewFetchCommand(c, f).Run())
	case Install:
		PrintOutput(cmd.NewInstallCommand(c, f, hc, i).Run())
	case Use:
		PrintOutput(cmd.NewUseCommand(c, f).Run())
	}
}

func createDefaultConfigIfIsNecessary(path string) error {
	var dir = filepath.Dir(path)

	if !file.Exists(path) {
		data, _ := json.Marshal(internal.AmvmConfig{Tools: internal.CreateDefaultConfig(dir)})
		if err := file.Write(path, data); err != nil {
			return fmt.Errorf("error creating default c file: %s", path)
		}
	}

	return nil
}

func loadConfiguration() (*internal.AmvmConfig, error) {
	mvmPath := env.Get("AMVM_HOME", internal.AmvmHomeDefaultDir)
	if err := os.MkdirAll(mvmPath, 0755); err != nil && !file.Exists(mvmPath) {
		return nil, err
	}

	configFilePath := filepath.Join(mvmPath, "config.json")
	if err := createDefaultConfigIfIsNecessary(configFilePath); err != nil {
		return nil, err
	}

	c, err := readConfig(configFilePath)
	if err != nil {
		return nil, err
	}

	dc := internal.CreateDefaultConfig(mvmPath)
	for _, tool := range internal.AvailableTools {
		tc := internal.Config{}

		tc.HomeDir = env.Get(fmt.Sprintf("AMVM_%s_HOME", strings.ToUpper(string(tool))), dc[tool].HomeDir)
		if err := os.MkdirAll(tc.HomeDir, 0755); err != nil && !file.Exists(tc.HomeDir) {
			return &internal.AmvmConfig{}, err
		}

		tc.CacheDir = env.Get(fmt.Sprintf("AMVM_%s_CACHE", strings.ToUpper(string(tool))), dc[tool].CacheDir)
		if err := os.MkdirAll(tc.CacheDir, 0755); err != nil && !file.Exists(tc.CacheDir) {
			return &internal.AmvmConfig{}, err
		}

		tc.VersionsDir = env.Get(fmt.Sprintf("AMVM_%s_VERSIONS", strings.ToUpper(string(tool))), dc[tool].VersionsDir)
		if err := os.MkdirAll(tc.VersionsDir, 0755); err != nil && !file.Exists(tc.VersionsDir) {
			return &internal.AmvmConfig{}, err
		}

		tc.CurrentDir = env.Get(fmt.Sprintf("AMVM_%s_VERSIONS", strings.ToUpper(string(tool))), dc[tool].CurrentDir)
		if err := os.MkdirAll(tc.CurrentDir, 0755); err != nil && !file.Exists(tc.CurrentDir) {
			return &internal.AmvmConfig{}, err
		}

		c.Tools[tool] = tc
	}

	if err := writeConfig(configFilePath, *c); err != nil {
		return nil, err
	}

	return c, nil
}

func readConfig(path string) (*internal.AmvmConfig, error) {
	data, err := file.Read(path)
	if err != nil {
		return &internal.AmvmConfig{}, fmt.Errorf("error reading c file: %s", path)
	}

	var tools = new(map[internal.Tool]internal.Config)
	c := &internal.AmvmConfig{HomeDir: filepath.Dir(path), Tools: *tools}
	err = json.Unmarshal(data, c)

	if err != nil {
		return &internal.AmvmConfig{}, fmt.Errorf("invalid c file: %s", path)
	}

	return c, nil
}

func writeConfig(path string, config internal.AmvmConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err = file.Write(path, data); err != nil {
		return err
	}

	return nil
}

func PrintOutput(output internal.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}
