package colors

import "github.com/fatih/color"

var (
	Error   = color.New(color.FgRed).SprintFunc()
	Info    = color.New(color.FgYellow).SprintFunc()
	Infof   = color.New(color.FgYellow).SprintfFunc()
	Success = color.New(color.FgGreen).SprintFunc()

	Mnemonic = color.New(color.FgHiBlue).SprintFunc()
	Modified = color.New(color.FgMagenta).SprintFunc()
	Name     = color.New(color.FgWhite, color.Bold).SprintFunc()
)

func SprintFunc(c int) func(...interface{}) string {
	return color.New(color.Attribute(c)).SprintFunc()
}
