package cmd

import (
	"fmt"
	"os"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/fetch"
	"github.com/romycode/amvm/internal/install"
	"github.com/romycode/amvm/pkg/http"
	"github.com/romycode/amvm/pkg/ui"
)

// InstallCommand command for download required version and save into AMVM_{TOOL}_versions
type InstallCommand struct {
	c  *internal.AmvmConfig
	f  *fetch.Fetcher
	hc http.Client
	i  *install.Installer
}

// NewInstallCommand return an instance of InstallCommand
func NewInstallCommand(c *internal.AmvmConfig, f *fetch.Fetcher, hc http.Client, i *install.Installer) *InstallCommand {
	return &InstallCommand{c, f, hc, i}
}

// Run get version and download `tar.gz` for save uncompressed into AMVM_{TOOL}_versions
func (i InstallCommand) Run() internal.Output {
	if len(os.Args[2:]) < 2 {
		return internal.NewOutput("invalid cmd, use: amvm install nodejs v17.3.0", ui.Green, 1)
	}

	tool := internal.Tool(os.Args[2])
	input := os.Args[3]

	versions, err := i.f.Run(tool)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)

	}

	version, err := versions.GetVersion(input)
	if err != nil {
		return internal.NewOutput(err.Error(), ui.Red, 1)

	}

	i.i.Run(tool, version)

	return internal.NewOutput(fmt.Sprintf("🔚 Download version: %s 🔚", input), ui.Green, 0)
}
