package cmd

import (
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/pnpm"
)

// InfoCommand command to get the latest versions of available tools
type InfoCommand struct {
	nf internal.Fetcher
	df internal.Fetcher
	pf internal.Fetcher
}

// NewInfoCommand returns new instance of InfoCommand
func NewInfoCommand(nf internal.Fetcher, df internal.Fetcher, pf internal.Fetcher) *InfoCommand {
	return &InfoCommand{
		nf: nf,
		df: df,
		pf: pf,
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

	pnpmVersions, err := i.pf.Run(pnpm.PnpmJs().Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	text := "Latest versions:	\n" +
		"\t- Node(latest): %s\n" +
		"\t- Node(lts)   : %s\n" +
		"\t- IoJs(latest): %s\n" +
		"\t- Deno(latest): %s\n" +
		"\t- Pnpm(latest): %s"

	return NewOutput(
		fmt.Sprintf(
			text,
			nodejsVersions.Latest().Semver(), nodejsVersions.Lts().Semver(),
			iojsVersions.Latest().Semver(), denoVersions.Latest().Semver(),
			pnpmVersions.Latest().Semver(),
		),
		1,
	)
}
