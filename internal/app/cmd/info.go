package cmd

import (
	"fmt"

	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
)

type InfoCommand struct {
	nf fetch.Fetcher
}

func NewInfoCommand(nf fetch.Fetcher) *InfoCommand {
	return &InfoCommand{
		nf: nf,
	}
}

func (i InfoCommand) Run() Output {
	iojsVersions, err := i.nf.Run(config.IoJs.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	nodejsVersions, err := i.nf.Run(config.NodeJs.Value())
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	return NewOutput(
		fmt.Sprintf(
			"Latest versions:\n  - Node(latest): %s\n  - Node(lts): %s\n  - IoJs(latest): %s",
			nodejsVersions.Latest().Semver(), nodejsVersions.Lts().Semver(), iojsVersions.Latest().Semver(),
		),
		1,
	)
}
