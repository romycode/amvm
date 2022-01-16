package main

import (
	"fmt"
	"os"

	"github.com/romycode/mvm/internal/app/cmd"
)

type Command string

const (
	Info Command = "info"
)

func main() {
	command := Command(os.Args[1])
	switch command {
	case Info:
		PrintOutput(cmd.NewInfoCommand().Run())
	}
}

func PrintOutput(output cmd.Output) {
	fmt.Println(output.Content)
	os.Exit(output.Code)
}
