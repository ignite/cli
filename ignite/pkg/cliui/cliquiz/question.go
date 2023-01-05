// Package cliquiz is a tool to collect answers from the users on cli.
package cliquiz

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/pflag"
)

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

func ask(q Question) error {
	var prompt survey.Prompt

	if !q.hidden {
		input := &survey.Input{
			Message: q.question,
		}
		if !q.required {
			input.Message += " (optional)"
		}
		if q.defaultAnswer != nil {
			input.Default = fmt.Sprintf("%v", q.defaultAnswer)
		}
		prompt = input
	} else {
		prompt = &survey.Password{
			Message: q.question,
		}
	}

	if err := survey.AskOne(prompt, q.answer); err != nil {
		return err
	}

	isValid := func() bool {
		if answer, ok := q.answer.(string); ok {
			if strings.TrimSpace(answer) == "" {
				return false
			}
		}
		if reflect.ValueOf(q.answer).Elem().IsZero() {
			return false
		}
		return true
	}

	if q.required && !isValid() {
		fmt.Println("This information is required, please retry:")

		if err := ask(q); err != nil {
			return err
		}
	}

	return nil
}

// Ask asks questions and collect answers.
func Ask(question ...Question) (err error) {
	defer func() {
		if errors.Is(err, terminal.InterruptErr) {
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
			return nil, fmt.Errorf("flag %q is not defined", f.Name)
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
