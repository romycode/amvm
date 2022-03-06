package cmd

import (
	"fmt"
	"os"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/pkg/color"
	"github.com/romycode/amvm/pkg/file"
)

// UseCommand command to set tool version active
type UseCommand struct {
	conf *config.MvmConfig
	nf   internal.Fetcher
}

// NewUseCommand returns an instance of UseCommand
func NewUseCommand(conf *config.MvmConfig, nf internal.Fetcher) *UseCommand {
	return &UseCommand{conf: conf, nf: nf}
}

// Run creates a symlink from tool version dir to MVM_{TOOL}_CURRENT
func (u UseCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm use nodejs v17.3.0", 1)
	}

	tool, err := node.NewFlavour(os.Args[2])
	if err != nil {
		return NewOutput(err.Error(), 1)
	}
	input := os.Args[3]

	if config.IoJsFlavour == tool || config.DefaultFlavour == tool {
		versions, err := u.nf.Run(tool.Value())
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		version, err := versions.GetVersion(input)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		if !file.Exists(u.conf.Node.VersionsDir + version.Semver()) {
			return NewOutput(
				color.Colorize(
					fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool.Value(), version.Semver()).Error(),
					color.Red,
				), 1)
		}

		_ = os.RemoveAll(u.conf.Node.CurrentDir)
		err = os.Symlink(u.conf.Node.VersionsDir+version.Semver(), u.conf.Node.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}
	}

	return NewOutput(color.Colorize(fmt.Sprintf("👌 Now 👉 version: %s", input), color.White), 1)
}
