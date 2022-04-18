package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/romycode/amvm/internal"
	"github.com/romycode/amvm/internal/config"
	"github.com/romycode/amvm/internal/deno"
	"github.com/romycode/amvm/internal/node"

	"github.com/romycode/amvm/pkg/color"
	"github.com/romycode/amvm/pkg/file"
)

// FetchCommand command for update tools versions and save into cache files
type FetchCommand struct {
	c  *config.AmvmConfig
	nf internal.Fetcher
	df internal.Fetcher
}

// NewFetchCommand returns new instance of FetchCommand
func NewFetchCommand(c *config.AmvmConfig, nf internal.Fetcher, df internal.Fetcher) *FetchCommand {
	return &FetchCommand{c: c, nf: nf, df: df}
}

// Run will execute fetch for every tool
func (f FetchCommand) Run() Output {
	cacheFile := fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", node.NodeJs().Value())
	versions, err := f.nf.Run(config.DefaultNodeJsFlavour.Value())
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

	cacheFile = fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", node.IoJs().Value())
	versions, err = f.nf.Run(config.IoJsFlavour.Value())
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

	cacheFile = fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", deno.DenoJs().Value())
	versions, err = f.df.Run(config.DefaultNodeJsFlavour.Value())
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
