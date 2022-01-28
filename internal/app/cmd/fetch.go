package cmd

import (
	"encoding/json"
	"fmt"
	config2 "github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
	"github.com/romycode/mvm/internal/node"

	"github.com/romycode/mvm/pkg/color"
	"github.com/romycode/mvm/pkg/file"
)

// FetchCommand command for update tools versions and save into cache files
type FetchCommand struct {
	c  *config2.MvmConfig
	nf fetch.Fetcher
}

// NewFetchCommand returns new instance of FetchCommand
func NewFetchCommand(c *config2.MvmConfig, nf fetch.Fetcher) *FetchCommand {
	return &FetchCommand{
		c:  c,
		nf: nf,
	}
}

// Run will execute fetch for every tool
func (f FetchCommand) Run() Output {
	cacheFile := fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", node.NodeJs)
	versions, err := f.nf.Run(node.DefaultFlavour.Value())
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	data, err := json.Marshal(versions)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	err = file.Write(cacheFile, data)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	cacheFile = fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", node.NodeJs)
	versions, err = f.nf.Run(node.IoJsFlavour.Value())
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	data, err = json.Marshal(versions)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	err = file.Write(cacheFile, data)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	return NewOutput(color.Colorize("➡ Update cache files ⬅", color.Blue), 0)
}
