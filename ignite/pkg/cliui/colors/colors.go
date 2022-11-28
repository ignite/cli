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
)

var (
	info     = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow)).Render
	infof    = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow)).Render
	error    = lipgloss.NewStyle().Foreground(lipgloss.Color(Red)).Render
	success  = lipgloss.NewStyle().Foreground(lipgloss.Color(Green)).Render
	modified = lipgloss.NewStyle().Foreground(lipgloss.Color(Magenta)).Render
)

func SprintFunc(code string) func(i ...interface{}) string {
	return func(i ...interface{}) string {
		render := lipgloss.NewStyle().Foreground(lipgloss.Color(code)).Render
		return render(fmt.Sprint(i...))
	}
}

func Info(i ...interface{}) string {
	return info(fmt.Sprint(i...))
}

func Infof(i ...interface{}) string {
	return infof(fmt.Sprint(i...))
}

func Error(i ...interface{}) string {
	return error(fmt.Sprint(i...))
}

func Success(i ...interface{}) string {
	return success(fmt.Sprint(i...))
}

func Modified(i ...interface{}) string {
	return modified(fmt.Sprint(i...))
}
