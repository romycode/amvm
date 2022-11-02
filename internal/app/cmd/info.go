package cmd

import (
	"fmt"
	"sync"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/fetch"
	"github.com/romycode/amvm/pkg/ui"
)

// InfoCommand command to get the latest versions of available tools
type InfoCommand struct {
	f *fetch.Fetcher
}

// NewInfoCommand returns new instance of InfoCommand
func NewInfoCommand(f *fetch.Fetcher) *InfoCommand {
	return &InfoCommand{f: f}
}

// Run fetch and print to stdout the latest versions
func (i InfoCommand) Run() Output {
	var wg sync.WaitGroup
	errorChan := make(chan error)

	var tools = map[internal.Tool]map[string]string{
		internal.Node: {"version": "", "name": "Node"},
		internal.Pnpm: {"version": "", "name": "Pnpm"},
		internal.Deno: {"version": "", "name": "Deno"},
		internal.Java: {"version": "", "name": "Java"},
	}

	for tool := range tools {
		tool := tool

		wg.Add(1)
		go func() {
			versions, err := i.f.Run(tool)
			if err != nil {
				errorChan <- err
			}

			tools[tool]["version"] = versions.Latest().Original()

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	if err := <-errorChan; err != nil {
		return NewOutput(err.Error(), ui.Red, 1)
	}

	output := "Latest versions:\n"
	for _, v := range tools {
		output += fmt.Sprintf("%s(latest): %s\n", v["name"], v["version"])
	}

	return NewOutput(output, ui.Green, 1)
}
