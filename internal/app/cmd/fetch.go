package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
	"github.com/romycode/mvm/pkg/color"
	"github.com/romycode/mvm/pkg/file"
)

type FetchCommand struct {
	c  *config.MvmConfig
	nf fetch.Fetcher
}

func NewFetchCommand(c *config.MvmConfig, nf fetch.Fetcher) *FetchCommand {
	return &FetchCommand{
		c:  c,
		nf: nf,
	}
}

func (f FetchCommand) Run() Output {
	cacheFile := fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", config.NodeJs)
	versions, err := f.nf.Run(config.NodeJs.Value())
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

	cacheFile = fmt.Sprintf(f.c.HomeDir+"/%s-versions.json", config.NodeJs)
	versions, err = f.nf.Run(config.IoJs.Value())
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
