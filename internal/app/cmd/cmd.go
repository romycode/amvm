package cmd

import "github.com/romycode/amvm/pkg/color"

type Output struct {
	Content string
	Code    int
}

func NewOutput(content string, messageColor color.Color, code int) Output {
	return Output{Content: color.Colorize(content, messageColor), Code: code}
}

type Command interface {
	Run() Output
}
