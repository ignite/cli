package colors

import "github.com/fatih/color"

var Info = color.New(color.FgYellow).SprintFunc()
var Error = color.New(color.FgRed).SprintFunc()
var Success = color.New(color.FgGreen).SprintFunc()
var Modified = color.New(color.FgMagenta).SprintFunc()

func SprintFunc(c int) func(...interface{}) string {
	return color.New(color.Attribute(c)).SprintFunc()
}
