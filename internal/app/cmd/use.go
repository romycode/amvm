package cmd

import (
	"fmt"
	"os"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/app/fetch"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/ui"
)

// UseCommand command to set tool version active
type UseCommand struct {
	conf *config.AmvmConfig
	ff   *fetch.Factory
}

// NewUseCommand returns an instance of UseCommand
func NewUseCommand(conf *config.AmvmConfig, ff *fetch.Factory) *UseCommand {
	return &UseCommand{conf: conf, ff: ff}
}

// Run creates a symlink from tool version dir to AMVM_{TOOL}_CURRENT
func (u UseCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm use nodejs v17.3.0", ui.Green, 1)
	}

	tool := os.Args[2]
	input := os.Args[3]

	vf, err := u.ff.Build(tool)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}
	versions, err := vf.Run(tool)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	version, err := versions.GetVersion(input)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	switch tool {
	case config.NodeFlavour.Value():
		return u.link(u.conf.Node.VersionsDir, u.conf.Node.CurrentDir, tool, version)
	case config.DenoFlavour.Value():
		return u.link(u.conf.Deno.VersionsDir, u.conf.Deno.CurrentDir, tool, version)
	case config.PnpmFlavour.Value():
		return u.link(u.conf.Pnpm.VersionsDir, u.conf.Pnpm.CurrentDir, tool, version)
	case config.JavaFlavour.Value():
		return u.link(u.conf.Java.VersionsDir, u.conf.Java.CurrentDir, tool, version)
	}

	return NewOutput(fmt.Sprintf("👌 Now 👉 version: %s", input), ui.White, 1)
}

func (u UseCommand) link(versionsDir string, currentDir string, tool string, version internal.Version) Output {
	if !file.Exists(versionsDir + version.SemverStr()) {
		return NewOutput(
			fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.SemverStr()).Error(),
			ui.Red,
			1,
		)
	}

	err := file.Link(versionsDir+version.SemverStr(), currentDir)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)
	}

	return Output{}
}
