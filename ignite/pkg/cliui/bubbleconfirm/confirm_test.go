package bubbleconfirm

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	m := NewModel("Continue?")

	require.Equal(t, "Continue?", m.Question)
	require.Equal(t, 0, m.cursor)
	require.Equal(t, Undecided, m.Choice())
}

func TestModelUpdateNavigationAndSelect(t *testing.T) {
	m := NewModel("Question")

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = next.(Model)
	require.Equal(t, 1, m.cursor)

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = next.(Model)
	require.Equal(t, No, m.Choice())
	require.NotNil(t, cmd)
}

func TestModelUpdateDirectYesNoChoices(t *testing.T) {
	m := NewModel("Question")

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	m = next.(Model)
	require.Equal(t, Yes, m.Choice())
	require.NotNil(t, cmd)

	m = NewModel("Question")
	next, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = next.(Model)
	require.Equal(t, No, m.Choice())
	require.NotNil(t, cmd)
}

func TestNewQuestionOptions(t *testing.T) {
	var answer string

	q := NewQuestion(
		"Name",
		&answer,
		DefaultAnswer("alice"),
		Required(),
		HideAnswer(),
		GetConfirmation(),
	)

	require.Equal(t, "Name", q.question)
	require.Equal(t, "alice", q.defaultAnswer)
	require.True(t, q.required)
	require.True(t, q.hidden)
	require.True(t, q.shouldConfirm)
	require.Equal(t, &answer, q.answer)
}

func TestInputModelUpdateRequiredValidation(t *testing.T) {
	m := inputModel{
		Question: "Name",
		Required: true,
	}

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = next.(inputModel)

	require.False(t, m.done)
	require.Equal(t, "this information is required", m.Error)
}

func TestInputModelUpdateTypingAndBackspace(t *testing.T) {
	m := inputModel{
		Question: "Name",
	}

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = next.(inputModel)
	require.Equal(t, "a", m.Value)
	require.Equal(t, 1, m.cursorPos)

	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = next.(inputModel)
	require.Equal(t, "", m.Value)
	require.Equal(t, 0, m.cursorPos)
}

func TestValuesFromFlagsOrAskUsesProvidedFlagValues(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("username", "", "username")
	fs.String("region", "", "region")
	require.NoError(t, fs.Set("username", "alice"))
	require.NoError(t, fs.Set("region", "earth"))

	values, err := ValuesFromFlagsOrAsk(
		fs,
		"",
		NewFlag("username", true),
		NewFlag("region", false),
	)
	require.NoError(t, err)
	require.Equal(t, "alice", values["username"])
	require.Equal(t, "earth", values["region"])
}

func TestValuesFromFlagsOrAskReturnsErrorForUndefinedFlag(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	_, err := ValuesFromFlagsOrAsk(fs, "", NewFlag("missing", true))
	require.Error(t, err)
}
