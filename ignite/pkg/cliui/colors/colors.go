package colors

import "github.com/fatih/color"

var (
	Info     = color.New(color.FgYellow).SprintFunc()
	Infof    = color.New(color.FgYellow).SprintfFunc()
	Error    = color.New(color.FgRed).SprintFunc()
	Success  = color.New(color.FgGreen).SprintFunc()
	Modified = color.New(color.FgMagenta).SprintFunc()
)

func SprintFunc(c int) func(...interface{}) string {
	return color.New(color.Attribute(c)).SprintFunc()
}
