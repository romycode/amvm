package cmd

import (
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/deno"
)

// InfoCommand command to get the latest versions of available tools
type InfoCommand struct {
	nf internal.Fetcher
	df internal.Fetcher
}

// NewInfoCommand returns new instance of InfoCommand
func NewInfoCommand(nf internal.Fetcher, df internal.Fetcher) *InfoCommand {
	return &InfoCommand{
		nf: nf,
		df: df,
	}
}

// Run fetch and print to stdout the latest versions
func (i InfoCommand) Run() Output {
	iojsVersions, err := i.nf.Run(config.IoJsFlavour.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	nodejsVersions, err := i.nf.Run(config.DefaultNodeJsFlavour.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	denoVersions, err := i.df.Run(deno.DenoJs().Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	return NewOutput(
		fmt.Sprintf(
			"Latest versions:\n  - NodeConfig(latest): %s\n  - NodeConfig(lts): %s\n  - IoJs(latest): %s\n  - DenoConfig(latest): %s",
			nodejsVersions.Latest().Semver(), nodejsVersions.Lts().Semver(), iojsVersions.Latest().Semver(), denoVersions.Latest().Semver(),
		),
		1,
	)
}
