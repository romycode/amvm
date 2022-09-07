package cmd

import (
	"fmt"
	"os"

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
		if !file.Exists(u.conf.Node.VersionsDir + version.SemverStr()) {
			return NewOutput(
				fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.SemverStr()).Error(),
				ui.Red,
				1,
			)
		}

		err = file.Link(u.conf.Node.VersionsDir+version.SemverStr(), u.conf.Node.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		break
	case config.DenoFlavour.Value():
		if !file.Exists(u.conf.Deno.VersionsDir + version.SemverStr()) {
			return NewOutput(
				fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.SemverStr()).Error(),
				ui.Red,
				1,
			)
		}

		err = file.Link(u.conf.Deno.VersionsDir+version.SemverStr(), u.conf.Deno.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)

		}

		break
	case config.PnpmFlavour.Value():
		if !file.Exists(u.conf.Pnpm.VersionsDir + version.SemverStr()) {
			return NewOutput(
				fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.SemverStr()).Error(),
				ui.Red,
				1)
		}

		err = file.Link(u.conf.Pnpm.VersionsDir+version.SemverStr(), u.conf.Pnpm.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)
		}

		break
	case config.JavaFlavour.Value():
		if !file.Exists(u.conf.Java.VersionsDir + version.SemverStr()) {
			return NewOutput(
				fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.SemverStr()).Error(),
				ui.Red,
				1)
		}

		err = file.Link(u.conf.Java.VersionsDir+version.SemverStr()+string(os.PathSeparator)+"Contents"+string(os.PathSeparator)+"Home", u.conf.Java.CurrentDir)
		if err != nil {
			return NewOutput(err.Error(), ui.Red, 1)
		}

		break
	}

	return NewOutput(fmt.Sprintf("👌 Now 👉 version: %s", input), ui.White, 1)
}
