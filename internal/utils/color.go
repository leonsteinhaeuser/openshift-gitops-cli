package utils

import "fmt"

type Color string

var (
	reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	Gray    Color = "\033[37m"
	White   Color = "\033[97m"
)

// Wrap wraps the given string with the color
func (c Color) Wrap(s string) string {
	return fmt.Sprintf("%s%s%s", c, s, reset)
}

// Colorize colorizes the given string with the given color
func Colorize(s string, c Color) string {
	return c.Wrap(s)
}
