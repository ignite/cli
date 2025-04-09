package bubbleconfirm

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// confirmation result values.
const (
	Undecided = iota
	Yes
	No
)

var (
	// styles for the confirmation dialog.
	questionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	yesStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	noStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

// Model represents the bubbletea model for a confirmation prompt.
type Model struct {
	Question string
	cursor   int
	choice   int
	done     bool
}

// NewModel creates a new confirmation model with the given question.
func NewModel(question string) Model {
	return Model{
		Question: question,
		cursor:   0, // 0 = yes, 1 = no
		choice:   Undecided,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { //nolint:gocritic // more readable than if-else
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.done = true
			m.choice = No
			return m, tea.Quit
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right", "l":
			if m.cursor < 1 {
				m.cursor++
			}
		case "enter", " ":
			// set choice based on cursor position
			m.done = true
			if m.cursor == 0 {
				m.choice = Yes
			} else {
				m.choice = No
			}
			return m, tea.Quit
		case "y", "Y":
			m.done = true
			m.choice = Yes
			return m, tea.Quit
		case "n", "N":
			m.done = true
			m.choice = No
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the confirmation prompt.
func (m Model) View() string {
	if m.done {
		return ""
	}

	question := questionStyle.Render(m.Question)
	yes := "Yes"
	no := "No"

	// apply styles based on cursor position
	if m.cursor == 0 {
		yes = cursorStyle.Render("[") + yesStyle.Render(yes) + cursorStyle.Render("]")
		no = "[ " + no + " ]"
	} else {
		yes = "[ " + yes + " ]"
		no = cursorStyle.Render("[") + noStyle.Render(no) + cursorStyle.Render("]")
	}

	return fmt.Sprintf("%s\n%s %s\n", question, yes, no)
}

// Choice returns the selected choice (Yes, No, or Undecided).
func (m Model) Choice() int {
	return m.choice
}
