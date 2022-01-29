package cmd

import (
	"fmt"

	"github.com/romycode/mvm/internal"
	"github.com/romycode/mvm/internal/config"
)

// InfoCommand command to get the latest versions of available tools
type InfoCommand struct {
	nf internal.Fetcher
}

// NewInfoCommand returns new instance of InfoCommand
func NewInfoCommand(nf internal.Fetcher) *InfoCommand {
	return &InfoCommand{
		nf: nf,
	}
}

// Run fetch and print to stdout the latest versions
func (i InfoCommand) Run() Output {
	iojsVersions, err := i.nf.Run(config.IoJsFlavour.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	nodejsVersions, err := i.nf.Run(config.DefaultFlavour.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	return NewOutput(
		fmt.Sprintf(
			"Latest versions:\n  - NodeConfig(latest): %s\n  - NodeConfig(lts): %s\n  - IoJs(latest): %s",
			nodejsVersions.Latest().Semver(), nodejsVersions.Lts().Semver(), iojsVersions.Latest().Semver(),
		),
		1,
	)
}
