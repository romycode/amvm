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

	output := "Latest versions:\n"
	for _, tool := range internal.AvailableTools {
		tool := tool

		wg.Add(1)
		go func() {
			defer wg.Done()

			versions, err := i.f.Run(tool)
			if err != nil {
				errorChan <- err
			}
			output += fmt.Sprintf("%s(latest): %s\n", string(tool), versions.Latest().Original())
		}()
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	if err := <-errorChan; err != nil {
		return NewOutput(err.Error(), ui.Red, 1)
	}

	return NewOutput(output, ui.Green, 1)
}
