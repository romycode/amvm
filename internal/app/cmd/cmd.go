package cmd

import "github.com/romycode/amvm/pkg/ui"

type Output struct {
	Content string
	Code    int
}

func NewOutput(content string, messageColor ui.Color, code int) Output {
	return Output{Content: ui.Colorize(content, messageColor), Code: code}
}

type Command interface {
	Run() Output
}
