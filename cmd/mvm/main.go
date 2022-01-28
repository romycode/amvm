package main

import (
	"fmt"
	httpstd "net/http"
	"os"

	"github.com/romycode/mvm/internal/app/cmd"
	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/internal/app/fetch"
	"github.com/romycode/mvm/pkg/color"
	"github.com/romycode/mvm/pkg/http"
)

type Command string

const (
	Info    Command = "info"
	Fetch   Command = "fetch"
	Install Command = "install"
	Use     Command = "use"
)

func main() {
	conf, err := config.LoadConfiguration()
	if err != nil {
		PrintOutput(cmd.NewOutput(color.Colorize(err.Error(), color.Red), 1))
	}
	if 1 == len(os.Args) {
		PrintOutput(cmd.NewOutput(color.Colorize("use: mvm <info|install|use|fetch> <nodejs> <flavour> <version>", color.White), 0))
	}

	nhc := http.NewClient(httpstd.DefaultClient, fetch.NodeJsURLTemplate)
	nf := fetch.NewNodeJsFetcher(nhc)

	command := Command(os.Args[1])
	switch command {
	case Info:
		PrintOutput(cmd.NewInfoCommand(nf).Run())
	case Fetch:
		PrintOutput(cmd.NewFetchCommand(conf, nf).Run())
	case Install:
		PrintOutput(cmd.NewInstallCommand(conf, nf).Run())
	case Use:
		PrintOutput(cmd.NewUseCommand(conf, nf).Run())
	}
}

func PrintOutput(output cmd.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}
