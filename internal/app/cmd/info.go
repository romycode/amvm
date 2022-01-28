package cmd

import (
	"fmt"
	"github.com/romycode/mvm/internal/app/fetch"
	"github.com/romycode/mvm/internal/node"
)

// InfoCommand command to get the latest versions of available tools
type InfoCommand struct {
	nf fetch.Fetcher
}

// NewInfoCommand returns new instance of InfoCommand
func NewInfoCommand(nf fetch.Fetcher) *InfoCommand {
	return &InfoCommand{
		nf: nf,
	}
}

// Run fetch and print to stdout the latest versions
func (i InfoCommand) Run() Output {
	iojsVersions, err := i.nf.Run(node.IoJsFlavour.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	nodejsVersions, err := i.nf.Run(node.DefaultFlavour.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	return NewOutput(
		fmt.Sprintf(
			"Latest versions:\n  - Config(latest): %s\n  - Config(lts): %s\n  - IoJs(latest): %s",
			nodejsVersions.Latest().Semver(), nodejsVersions.Lts().Semver(), iojsVersions.Latest().Semver(),
		),
		1,
	)
}
