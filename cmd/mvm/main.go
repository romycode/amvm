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
	Info  Command = "info"
	Fetch Command = "fetch"
)

func main() {
	conf, err := config.LoadConfiguration()
	if err != nil {
		PrintOutput(cmd.NewOutput(color.Colorize(err.Error(), color.Red), 1))
	}

	command := Command(os.Args[1])
	switch command {
	case Info:
		PrintOutput(cmd.NewInfoCommand().Run())
	case Fetch:
		PrintOutput(cmd.NewFetchCommand(*conf).Run())
	}
}

func PrintOutput(output cmd.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}
