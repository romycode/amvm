package cmd

import (
	"encoding/json"
	"path/filepath"
	"sync"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/fetch"
	"github.com/romycode/amvm/pkg/file"
	"github.com/romycode/amvm/pkg/ui"
)

// FetchCommand command for update tools versions and save into cache files
type FetchCommand struct {
	c *internal.AmvmConfig
	f *fetch.Fetcher
}

// NewFetchCommand returns new instance of FetchCommand
func NewFetchCommand(c *internal.AmvmConfig, f *fetch.Fetcher) *FetchCommand {
	return &FetchCommand{c, f}
}

// Run will execute fetch for every tool
func (r FetchCommand) Run() Output {
	spinner := ui.NewSpinner("Fetching versions ...")
	spinner.Start()
	defer spinner.Stop()

	var wg sync.WaitGroup
	errorChan := make(chan error)

	for _, tool := range internal.AvailableTools {
		tool := tool

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := r.createCacheFile(filepath.Join(r.c.HomeDir, string(tool)+"-versions.json"), tool); err != nil {
				errorChan <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	err := <-errorChan
	if err != nil {
		return NewOutput(err.Error(), ui.Red, 1)
	}

	return NewOutput("➡ Update cache files ⬅", ui.Blue, 0)
}

func (r FetchCommand) createCacheFile(filename string, tool internal.Tool) error {
	versions, err := r.f.Run(tool)
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
