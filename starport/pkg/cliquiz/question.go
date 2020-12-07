// Package cliquiz is a tool to collect answers from the users on cli.
package cliquiz

import (
	"fmt"
	"reflect"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/pflag"
)

// Question holds information on what to ask to user and where
// the answer stored at.
type Question struct {
	question      string
	defaultAnswer interface{}
	answer        interface{}
	hidden        bool
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

// Ask asks questions and collect answers.
func Ask(question ...Question) error {
	for _, q := range question {

		var prompt survey.Prompt
		if !q.hidden {
			input := &survey.Input{
				Message: q.question,
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

		if q.required && reflect.ValueOf(q.answer).Elem().IsZero() {
			fmt.Println("This information is required, please retry:")
			if err := Ask(q); err != nil {
				return err
			}
		}
	}
	return nil
}

// ValuesFromFlagsOrAsk returns values of requiredFlags within map[string]string where map's
// key is the name of the flag and value is flag's value.
// when provided, values are collected through command otherwise they're asked to user by prompting.
// title used as a message while prompting.
func ValuesFromFlagsOrAsk(flags *pflag.FlagSet, title string, requiredFlags []string) (values map[string]string, err error) {
	values = make(map[string]string)

	answers := make(map[string]*string)
	var questions []Question

	for _, name := range requiredFlags {
		flag := flags.Lookup(name)
		if flag == nil {
			return nil, fmt.Errorf("flag %q is not defined", name)
		}
		if value, _ := flags.GetString(name); value != "" {
			values[name] = value
			continue
		}

		var value string
		answers[name] = &value
		questions = append(questions, NewQuestion(flag.Usage, &value, Required()))
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
