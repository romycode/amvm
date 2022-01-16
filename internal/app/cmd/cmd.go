package cmd

type Output struct {
	Content string
	Code    int
}

func NewOutput(content string, code int) Output {
	return Output{Content: content, Code: code}
}

type Command interface {
	Run() Output
}
