package internal

import "github.com/romycode/amvm/pkg/ui"

type Output struct {
	Content string
	Code    int
}

func NewOutput(content string, messageColor ui.Color, code int) Output {
	newMsg := content

	for i := 0; i < 50-len(content); i++ {
		newMsg += " "
	}

	return Output{Content: ui.Colorize(newMsg, messageColor), Code: code}
}

type Command interface {
	Run() Output
}
