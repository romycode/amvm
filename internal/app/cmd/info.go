package cmd

import (
	"fmt"
	"sync"

	"github.com/romycode/amvm/internal/app/fetch"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/internal/pnpm"
	"github.com/romycode/amvm/pkg/color"
)

// InfoCommand command to get the latest versions of available tools
type InfoCommand struct {
	ff *fetch.Factory
}

// NewInfoCommand returns new instance of InfoCommand
func NewInfoCommand(ff *fetch.Factory) *InfoCommand {
	return &InfoCommand{ff: ff}
}

// Run fetch and print to stdout the latest versions
func (i InfoCommand) Run() Output {
	var wg sync.WaitGroup
	errorChan := make(chan error)

	var tools = map[string]map[string]string{
		node.NodeJs().Value(): {"version": "", "name": "Node"},
		deno.DenoJs().Value(): {"version": "", "name": "Deno"},
		pnpm.PnpmJs().Value(): {"version": "", "name": "Pnpm"},
	}

	for k, _ := range tools {
		wg.Add(1)

		go func(tool string) {
			fetcher, err := i.ff.Build(tool)
			if err != nil {
				errorChan <- err
			}

			versions, err := fetcher.Run(tool)
			if err != nil {
				errorChan <- err
			}

			tools[tool]["version"] = versions.Latest().Original()

			wg.Done()
		}(k)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	if err := <-errorChan; err != nil {
		return NewOutput(err.Error(), color.Red, 1)
	}

	output := "Latest versions:\n"
	for _, v := range tools {
		output += fmt.Sprintf("%s(latest): %s\n", v["name"], v["version"])
	}

	return NewOutput(output, color.Green, 1)
}
