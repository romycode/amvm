package cmd

import (
	"fmt"

	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
)

type InfoCommand struct {
}

func NewInfoCommand() *InfoCommand {
	return &InfoCommand{}
}

func (i InfoCommand) Run() Output {
	iojsVersions, err := fetch.NodeVersions(config.IoJs)
	if err != nil {
		return NewOutput(err.Error(), 1)
	}

	nodejsVersions, err := fetch.NodeVersions(config.NodeJs)
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
