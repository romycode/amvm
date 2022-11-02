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
func (u UseCommand) Run() Output {
	if len(os.Args[2:]) < 2 {
		return NewOutput("invalid cmd, use: amvm use nodejs v17.3.0", ui.Green, 1)
	}

	tool := internal.Tool(os.Args[2])
	input := os.Args[3]

	vs, err := u.f.Run(tool)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	v, err := vs.GetVersion(input)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)

	}

	switch tool {
	case internal.Node:
		return u.link(u.c.Tools[internal.Node].VersionsDir, u.c.Tools[internal.Node].CurrentDir, string(tool), v)
	case internal.Deno:
		return u.link(u.c.Tools[internal.Deno].VersionsDir, u.c.Tools[internal.Deno].CurrentDir, string(tool), v)
	case internal.Pnpm:
		return u.link(u.c.Tools[internal.Pnpm].VersionsDir, u.c.Tools[internal.Pnpm].CurrentDir, string(tool), v)
	case internal.Java:
		return u.link(u.c.Tools[internal.Java].VersionsDir, u.c.Tools[internal.Java].CurrentDir, string(tool), v)
	}

	return NewOutput(fmt.Sprintf("👌 Now 👉 v: %s", input), ui.White, 1)
}

func (u UseCommand) link(versionsDir string, currentDir string, tool string, version version.Version) Output {
	if !file.Exists(filepath.Join(versionsDir, version.SemverStr())) {
		return NewOutput(
			fmt.Errorf("version not downloaded, install with: amvm install %s %s", tool, version.SemverStr()).Error(),
			ui.Red,
			1,
		)
	}

	err := file.Link(filepath.Join(versionsDir, version.SemverStr()), currentDir)
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)
	}

	return Output{}
}
