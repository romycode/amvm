package cmd

import (
	"fmt"
	"os"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/pkg/color"
	"github.com/romycode/amvm/pkg/file"
)

// UseCommand command to set tool version active
type UseCommand struct {
	conf *config.AmvmConfig
	nf   internal.Fetcher
	df   internal.Fetcher
}

// NewUseCommand returns an instance of UseCommand
func NewUseCommand(conf *config.AmvmConfig, nf internal.Fetcher, df internal.Fetcher) *UseCommand {
	return &UseCommand{conf: conf, nf: nf, df: df}
}

// Run creates a symlink from tool version dir to AMVM_{TOOL}_CURRENT
func (u UseCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm use nodejs v17.3.0", 1)
	}

	tool := os.Args[2]
	_, notNodeTool := node.NewFlavour(tool)
	_, notDenoTool := deno.NewFlavour(tool)
	if notNodeTool != nil && notDenoTool != nil {
		message := notNodeTool.Error()
		if notDenoTool != nil {
			message = notDenoTool.Error()
		}
		return NewOutput(message, 1)
	}

	input := os.Args[3]
	if config.IoJsFlavour.Value() == tool || config.DefaultNodeJsFlavour.Value() == tool {
		versions, err := u.nf.Run(tool)
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
					fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.Semver()).Error(),
					color.Red,
				), 1)
		}

		_ = os.RemoveAll(u.conf.Node.CurrentDir)
		err = os.Symlink(u.conf.Node.VersionsDir+version.Semver(), u.conf.Node.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}
	}

	if config.DefaultDenoJsFlavour.Value() == tool {
		versions, err := u.df.Run(tool)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		version, err := versions.GetVersion(input)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		if !file.Exists(u.conf.Deno.VersionsDir + version.Semver()) {
			return NewOutput(
				color.Colorize(
					fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.Semver()).Error(),
					color.Red,
				), 1)
		}

		_ = os.RemoveAll(u.conf.Deno.CurrentDir)
		err = os.Symlink(u.conf.Deno.VersionsDir+version.Semver(), u.conf.Deno.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}
	}

	return NewOutput(color.Colorize(fmt.Sprintf("👌 Now 👉 version: %s", input), color.White), 1)
}
