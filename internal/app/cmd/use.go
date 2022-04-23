package cmd

import (
	"fmt"
	"os"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/pkg/color"
	"github.com/romycode/amvm/pkg/file"
)

// UseCommand command to set tool version active
type UseCommand struct {
	conf *config.AmvmConfig
	nf   internal.Fetcher
	df   internal.Fetcher
	pf   internal.Fetcher
}

// NewUseCommand returns an instance of UseCommand
func NewUseCommand(conf *config.AmvmConfig, nf internal.Fetcher, df internal.Fetcher, pf internal.Fetcher) *UseCommand {
	return &UseCommand{conf: conf, nf: nf, df: df, pf: pf}
}

// Run creates a symlink from tool version dir to AMVM_{TOOL}_CURRENT
func (u UseCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm use nodejs v17.3.0", 1)
	}

	tool := os.Args[2]
	input := os.Args[3]
	switch tool {
	case config.IoJsFlavour.Value():
	case config.NodeJsFlavour.Value():
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

		break
	case config.DenoJsFlavour.Value():
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

		break
	case config.PnpmJsFlavour.Value():
		versions, err := u.pf.Run(tool)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		version, err := versions.GetVersion(input)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		if !file.Exists(u.conf.Pnpm.VersionsDir + version.Semver()) {
			return NewOutput(
				color.Colorize(
					fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.Semver()).Error(),
					color.Red,
				), 1)
		}

		_ = os.RemoveAll(u.conf.Pnpm.CurrentDir)
		err = os.Symlink(u.conf.Pnpm.VersionsDir+version.Semver(), u.conf.Pnpm.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), 1)
		}

		break
	}

	return NewOutput(color.Colorize(fmt.Sprintf("👌 Now 👉 version: %s", input), color.White), 1)
}
