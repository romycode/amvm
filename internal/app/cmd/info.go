package cmd

import (
	"fmt"
	"sync"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/app/fetch"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/deno"
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

	var iojsVersions internal.Versions
	var nodejsVersions internal.Versions
	var denoVersions internal.Versions
	var pnpmVersions internal.Versions

	wg.Add(4)

	go func() {
		var nodejsFetcher internal.Fetcher
		nodejsFetcher, err := i.ff.Build(config.NodeJsFlavour.Value())
		if err != nil {
			errorChan <- err
		}
		nodejsVersions, err = nodejsFetcher.Run(config.NodeJsFlavour.Value())
		if err != nil {
			errorChan <- err
		}

		wg.Done()
	}()

	go func() {
		var iojsFetcher internal.Fetcher
		iojsFetcher, err := i.ff.Build(config.IoJsFlavour.Value())
		if err != nil {
			errorChan <- err
		}
		iojsVersions, err = iojsFetcher.Run(config.IoJsFlavour.Value())
		if err != nil {
			errorChan <- err
		}

		wg.Done()
	}()

	go func() {
		var denoFetcher internal.Fetcher
		denoFetcher, err := i.ff.Build(deno.DenoJs().Value())
		if err != nil {
			errorChan <- err
		}
		denoVersions, err = denoFetcher.Run(deno.DenoJs().Value())
		if err != nil {
			errorChan <- err
		}

		wg.Done()
	}()

	go func() {
		var pnpmFetcher internal.Fetcher
		pnpmFetcher, err := i.ff.Build(pnpm.PnpmJs().Value())
		if err != nil {
			errorChan <- err
		}
		pnpmVersions, err = pnpmFetcher.Run(pnpm.PnpmJs().Value())
		if err != nil {
			errorChan <- err
		}

		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	if err := <-errorChan; err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	text := "Latest versions:\n" +
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
