package color

import "fmt"

func Colorize(msg string, color Color) string {
	return fmt.Sprintf("%s%s%s", color, msg, Reset)
}
