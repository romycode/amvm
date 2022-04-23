package cmd

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/romycode/amvm/internal/app/fetch"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/node"
	"github.com/romycode/amvm/internal/pnpm"
	"github.com/romycode/amvm/pkg/color"
	"github.com/romycode/amvm/pkg/file"
)

// FetchCommand command for update tools versions and save into cache files
type FetchCommand struct {
	c  *config.AmvmConfig
	ff *fetch.Factory
}

// NewFetchCommand returns new instance of FetchCommand
func NewFetchCommand(c *config.AmvmConfig, ff *fetch.Factory) *FetchCommand {
	return &FetchCommand{c: c, ff: ff}
}

// Run will execute fetch for every tool
func (f FetchCommand) Run() Output {
	var tools = map[string]string{
		node.NodeJs().Value(): fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", node.NodeJs().Value()),
		node.IoJs().Value():   fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", node.IoJs().Value()),
		deno.DenoJs().Value(): fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", deno.DenoJs().Value()),
		pnpm.PnpmJs().Value(): fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", pnpm.PnpmJs().Value()),
	}

	var wg sync.WaitGroup
	errorChan := make(chan error)

	for k, v := range tools {
		wg.Add(1)

		go func(filename, tool string) {
			err := f.createCacheFile(filename, tool)
			if err != nil {
				errorChan <- err
			}

			wg.Done()
		}(k, v)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	err := <-errorChan
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	return NewOutput(color.Colorize("➡ Update cache files ⬅", color.Blue), 0)
}

func (f FetchCommand) createCacheFile(filename, tool string) error {
	fetcher, err := f.ff.Build(tool)
	if err != nil {
		return err
	}

	versions, err := fetcher.Run(tool)
	data, err := json.Marshal(versions)
	if err != nil {
		return err
	}

	err = file.Write(filename, data)
	if err != nil {
		return err
	}

	return nil
}
