package cmd

import (
	"fmt"
	"os"

	"github.com/romycode/mvm/internal"
	"github.com/romycode/mvm/internal/config"
	"github.com/romycode/mvm/internal/node"
	"github.com/romycode/mvm/pkg/color"
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

// Run creates a symlink from tool verison dir to MVM_{TOOL}_CURRENT
func (u UseCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		fmt.Println("invalid cmd, use: mvm use nodejs v17.3.0")
		os.Exit(1)
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

		err = os.RemoveAll(u.conf.Node.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		err = os.Symlink(u.conf.Node.VersionsDir+version.Semver(), u.conf.Node.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}
	}

	return NewOutput(color.Colorize(fmt.Sprintf("👌 Now 👉 version: %s", input), color.White), 1)
}
