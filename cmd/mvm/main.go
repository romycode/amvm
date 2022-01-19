package main

import (
	"fmt"
	"os"

	"github.com/romycode/mvm/internal/app/cmd"
	"github.com/romycode/mvm/internal/app/config"
	"github.com/romycode/mvm/pkg/color"
)

type Command string

const (
	Info    Command = "info"
	Fetch   Command = "fetch"
	Install Command = "install"
)

func main() {
	conf, err := config.LoadConfiguration()
	if err != nil {
		PrintOutput(cmd.NewOutput(color.Colorize(err.Error(), color.Red), 1))
	}
	if 1 == len(os.Args) {
		PrintOutput(cmd.NewOutput(color.Colorize("use: mvm <info|install|use|fetch> <nodejs> <flavour> <version>", color.White), 0))
	}

	command := Command(os.Args[1])
	switch command {
	case Info:
		PrintOutput(cmd.NewInfoCommand().Run())
	case Fetch:
		PrintOutput(cmd.NewFetchCommand(*conf).Run())
	case Install:
		PrintOutput(cmd.NewInstallCommand(*conf).Run())
	}
}

func PrintOutput(output cmd.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}
