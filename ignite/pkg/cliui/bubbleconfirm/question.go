package bubbleconfirm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/pflag"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

var (
	// styles for the question input.
	activeStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	promptStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

// ErrInterrupted is returned when the input process is interrupted.
var ErrInterrupted = errors.New("interrupted")

// ErrConfirmationFailed is returned when second answer is not the same with first one.
var ErrConfirmationFailed = errors.New("failed to confirm, your answers were different")

// Question holds information on what to ask user and where
// the answer stored at.
type Question struct {
	question      string
	defaultAnswer interface{}
	answer        interface{}
	hidden        bool
	shouldConfirm bool
	required      bool
}

// Option configures Question.
type Option func(*Question)

// DefaultAnswer sets a default answer to Question.
func DefaultAnswer(answer interface{}) Option {
	return func(q *Question) {
		q.defaultAnswer = answer
	}
}

// Required marks the answer as required.
func Required() Option {
	return func(q *Question) {
		q.required = true
	}
}

// HideAnswer hides the answer to prevent secret information being leaked.
func HideAnswer() Option {
	return func(q *Question) {
		q.hidden = true
	}
}

// GetConfirmation prompts confirmation for the given answer.
func GetConfirmation() Option {
	return func(q *Question) {
		q.shouldConfirm = true
	}
}

// NewQuestion creates a new question.
func NewQuestion(question string, answer interface{}, options ...Option) Question {
	q := Question{
		question: question,
		answer:   answer,
	}
	for _, o := range options {
		o(&q)
	}
	return q
}

// inputModel represents the bubbletea model for an input prompt.
type inputModel struct {
	Question     string
	Value        string
	Hidden       bool
	Required     bool
	DefaultValue string
	Error        string
	cursorPos    int
	done         bool
}

// Init initializes the input model.
func (m inputModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the input model.
func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { //nolint:gocritic // more readable than if-else
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.done = true
			return m, tea.Quit
		case "enter":
			// validate if input is required
			if m.Required && strings.TrimSpace(m.Value) == "" {
				m.Error = "this information is required"
				return m, nil
			}

			m.done = true
			return m, tea.Quit
		case "backspace":
			if m.cursorPos > 0 {
				m.Value = m.Value[:m.cursorPos-1] + m.Value[m.cursorPos:]
				m.cursorPos--
			}
		case "left":
			if m.cursorPos > 0 {
				m.cursorPos--
			}
		case "right":
			if m.cursorPos < len(m.Value) {
				m.cursorPos++
			}
		case "home":
			m.cursorPos = 0
		case "end":
			m.cursorPos = len(m.Value)
		default:
			// only accept printable characters
			if len(msg.Runes) == 1 {
				m.Value = m.Value[:m.cursorPos] + string(msg.Runes) + m.Value[m.cursorPos:]
				m.cursorPos++
				m.Error = ""
			}
		}
	}
	return m, nil
}

// View renders the input prompt.
func (m inputModel) View() string {
	if m.done {
		return ""
	}

	question := m.Question
	if !m.Required {
		question += " (optional)"
	}
	question = questionStyle.Render(question)

	var input string
	if m.Hidden {
		// show asterisks for hidden input
		input = strings.Repeat("*", len(m.Value))
	} else {
		input = m.Value
	}

	// show cursor position
	var display string
	if m.Value == "" && m.DefaultValue != "" { //nolint:gocritic // more readable than switch
		// show default value as placeholder
		display = placeholderStyle.Render(m.DefaultValue)
	} else if m.cursorPos < len(input) {
		display = input[:m.cursorPos] + activeStyle.Render(string(input[m.cursorPos])) + input[m.cursorPos+1:]
	} else {
		display = input + activeStyle.Render("_")
	}

	prompt := fmt.Sprintf("%s\n%s ", question, promptStyle.Render("›"))

	if m.Error != "" {
		return prompt + display + "\n" + errorStyle.Render(m.Error)
	}

	return prompt + display
}

func ask(q Question) error {
	// prepare default value as string
	defaultValue := ""
	if q.defaultAnswer != nil {
		defaultValue = fmt.Sprintf("%v", q.defaultAnswer)
	}

	// create and init the model
	m := inputModel{
		Question:     q.question,
		Hidden:       q.hidden,
		Required:     q.required,
		DefaultValue: defaultValue,
	}

	// run the bubbletea program
	p := tea.NewProgram(&m)
	result, err := p.Run()
	if err != nil {
		return err
	}

	finalModel := result.(inputModel)
	if !finalModel.done {
		return ErrInterrupted
	}

	// if empty and we have a default, use the default
	value := finalModel.Value
	if value == "" && defaultValue != "" {
		value = defaultValue
	}

	// convert the string value to the target type and store it
	switch ptr := q.answer.(type) {
	case *string:
		*ptr = value
	case *int:
		var i int
		_, err := fmt.Sscanf(value, "%d", &i)
		if err == nil {
			*ptr = i
		}
	case *float64:
		var f float64
		_, err := fmt.Sscanf(value, "%f", &f)
		if err == nil {
			*ptr = f
		}
	case *bool:
		*ptr = strings.ToLower(value) == "true" || value == "1" || strings.ToLower(value) == "yes" || strings.ToLower(value) == "y"
	default:
		// use reflection for other types
		v := reflect.ValueOf(ptr).Elem()
		if v.Kind() == reflect.String {
			v.SetString(value)
		}
	}

	return nil
}

// Ask asks questions and collect answers.
func Ask(question ...Question) (err error) {
	defer func() {
		if errors.Is(err, ErrInterrupted) {
			err = context.Canceled
		}
	}()

	for _, q := range question {
		if err := ask(q); err != nil {
			return err
		}

		if q.shouldConfirm {
			var secondAnswer string

			var options []Option
			if q.required {
				options = append(options, Required())
			}
			if q.hidden {
				options = append(options, HideAnswer())
			}
			if err := ask(NewQuestion("Confirm "+q.question, &secondAnswer, options...)); err != nil {
				return err
			}

			t := reflect.TypeOf(secondAnswer)
			compAnswer := reflect.ValueOf(q.answer).Elem().Convert(t).String()
			if secondAnswer != compAnswer {
				return ErrConfirmationFailed
			}
		}
	}
	return nil
}

// Flag represents a cmd flag.
type Flag struct {
	Name       string
	IsRequired bool
}

// NewFlag creates a new flag.
func NewFlag(name string, isRequired bool) Flag {
	return Flag{name, isRequired}
}

// ValuesFromFlagsOrAsk returns values of flags within map[string]string where map's
// key is the name of the flag and value is flag's value.
// when provided, values are collected through command otherwise they're asked by prompting user.
// title used as a message while prompting.
func ValuesFromFlagsOrAsk(fset *pflag.FlagSet, title string, flags ...Flag) (values map[string]string, err error) {
	values = make(map[string]string)

	answers := make(map[string]*string)
	var questions []Question

	for _, f := range flags {
		flag := fset.Lookup(f.Name)
		if flag == nil {
			return nil, errors.Errorf("flag %q is not defined", f.Name)
		}
		if value, _ := fset.GetString(f.Name); value != "" {
			values[f.Name] = value
			continue
		}

		var value string
		answers[f.Name] = &value

		var options []Option
		if f.IsRequired {
			options = append(options, Required())
		}
		questions = append(questions, NewQuestion(flag.Usage, &value, options...))
	}

	if len(questions) > 0 && title != "" {
		fmt.Println(title)
	}
	if err := Ask(questions...); err != nil {
		return values, err
	}

	for name, answer := range answers {
		values[name] = *answer
	}

	return values, nil
}
