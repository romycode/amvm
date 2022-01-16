package cmd

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
	"github.com/romycode/mvm/pkg/color"
	"github.com/romycode/mvm/pkg/file"
)

type FetchCommand struct {
	conf config.MvmConfig
}

func NewFetchCommand(conf config.MvmConfig) *FetchCommand {
	return &FetchCommand{
		conf: conf,
	}
}

func (f FetchCommand) Run() Output {
	cacheFile := fmt.Sprintf(f.conf.HomeDir+"/%s-versions.json", config.NodeJs)

	res, err := http.Get(fetch.NodeJsURL + fetch.NodeJsVersionsURL)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}
	err = file.Write(cacheFile, data)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	cacheFile = filepath.Join(f.conf.HomeDir, fmt.Sprintf("%s-versions.json", config.IoJs))
	res, err = http.Get(fetch.IoJsURL + fetch.IoJsVersionsURL)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}
	err = file.Write(cacheFile, data)
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	err = res.Body.Close()
	if err != nil {
		return NewOutput(color.Colorize(err.Error(), color.Red), 1)
	}

	return NewOutput(color.Colorize("➡ Update cache files ⬅", color.Blue), 0)
}
