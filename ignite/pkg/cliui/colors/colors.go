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
	Blue    = "#0a2fc4"
)

var (
	info     = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow))
	infof    = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow))
	err      = lipgloss.NewStyle().Foreground(lipgloss.Color(Red))
	success  = lipgloss.NewStyle().Foreground(lipgloss.Color(Green))
	modified = lipgloss.NewStyle().Foreground(lipgloss.Color(Magenta))
	name     = lipgloss.NewStyle().Bold(true)
	mnemonic = lipgloss.NewStyle().Foreground(lipgloss.Color(HiBlue))
	spinner  = lipgloss.NewStyle().Foreground(lipgloss.Color(Blue))
	faint    = lipgloss.NewStyle().Faint(true)
)

// SprintFunc returns a function to apply a foreground color to any number of texts.
// The returned function receives strings as arguments with the text that should be colorized.
// Color specifies a color by hex or ANSI value.
func SprintFunc(color string) func(i ...any) string {
	return func(i ...any) string {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		return style.Render(fmt.Sprint(i...))
	}
}

func Info(i ...any) string {
	return info.Render(fmt.Sprint(i...))
}

func Infof(format string, i ...any) string {
	return infof.Render(fmt.Sprintf(format, i...))
}

func Error(i ...any) string {
	return err.Render(fmt.Sprint(i...))
}

func Success(i ...any) string {
	return success.Render(fmt.Sprint(i...))
}

func Modified(i ...any) string {
	return modified.Render(fmt.Sprint(i...))
}

func Name(i ...any) string {
	return name.Render(fmt.Sprint(i...))
}

func Mnemonic(i ...any) string {
	return mnemonic.Render(fmt.Sprint(i...))
}

func Spinner(i ...any) string {
	return spinner.Render(fmt.Sprint(i...))
}

// Faint styles the text using a dimmer shade for the foreground color.
func Faint(i ...any) string {
	return faint.Render(fmt.Sprint(i...))
}
