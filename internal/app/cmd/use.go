package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/fetch"
	"github.com/romycode/amvm/internal/version"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/ui"
)

// UseCommand command to set tool version active
type UseCommand struct {
	c *internal.AmvmConfig
	f *fetch.Fetcher
}

// NewUseCommand returns an instance of UseCommand
func NewUseCommand(c *internal.AmvmConfig, f *fetch.Fetcher) *UseCommand {
	return &UseCommand{c: c, f: f}
}

// Run creates a symlink from tool version dir to AMVM_{TOOL}_CURRENT
func (u UseCommand) Run() internal.Output {
	if len(os.Args[2:]) < 2 {
		return internal.NewOutput("invalid cmd, use: amvm use nodejs v17.3.0", ui.Green, 1)
	}

	tool := internal.Tool(os.Args[2])
	input := os.Args[3]

	vs, err := u.f.Run(tool)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)

	}

	v, err := vs.GetVersion(input)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)

	}

	c := u.c.Tools[tool]
	if !file.Exists(filepath.Join(c.VersionsDir, v.SemverStr())) {
		return internal.NewOutput(
			fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, v.SemverStr()).Error(),
			ui.Red,
			1,
		)
	}

	if err := u.link(c.VersionsDir, c.CurrentDir, v); err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)
	}
	return internal.NewOutput(fmt.Sprintf("👌 Now 👉 v: %s", input), ui.White, 1)
}

func (u UseCommand) link(versionsDir string, currentDir string, version version.Version) error {
	err := file.Link(filepath.Join(versionsDir, version.SemverStr()), currentDir)
	if err != nil {
		return err
	}

	return nil
}
