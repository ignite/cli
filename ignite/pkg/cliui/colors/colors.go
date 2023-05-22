package colors

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	Yellow  = "#c4a000"
	Red     = "#ef2929"
	Green   = "#4e9a06"
	Magenta = "#75507b"
	Cyan    = "#34e2e2"
	White   = "#FFFFFF"
	HiBlue  = "#729FCF"
)

var (
	info     = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow))
	infof    = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow))
	err      = lipgloss.NewStyle().Foreground(lipgloss.Color(Red))
	success  = lipgloss.NewStyle().Foreground(lipgloss.Color(Green))
	modified = lipgloss.NewStyle().Foreground(lipgloss.Color(Magenta))
	name     = lipgloss.NewStyle().Bold(true)
	mnemonic = lipgloss.NewStyle().Foreground(lipgloss.Color(HiBlue))
	faint    = lipgloss.NewStyle().Faint(true)
)

// SprintFunc returns a function to apply a foreground color to any number of texts.
// The returned function receives strings as arguments with the text that should be colorized.
// Color specifies a color by hex or ANSI value.
func SprintFunc(color string) func(i ...interface{}) string {
	return func(i ...interface{}) string {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		return style.Render(fmt.Sprint(i...))
	}
}

func Info(i ...interface{}) string {
	return info.Render(fmt.Sprint(i...))
}

func Infof(format string, i ...interface{}) string {
	return infof.Render(fmt.Sprintf(format, i...))
}

func Error(i ...interface{}) string {
	return err.Render(fmt.Sprint(i...))
}

func Success(i ...interface{}) string {
	return success.Render(fmt.Sprint(i...))
}

func Modified(i ...interface{}) string {
	return modified.Render(fmt.Sprint(i...))
}

func Name(i ...interface{}) string {
	return name.Render(fmt.Sprint(i...))
}

func Mnemonic(i ...interface{}) string {
	return mnemonic.Render(fmt.Sprint(i...))
}

// Faint styles the text using a dimmer shade for the foreground color.
func Faint(i ...interface{}) string {
	return faint.Render(fmt.Sprint(i...))
}
