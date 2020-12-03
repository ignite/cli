// Package cliquiz is a tool to collect answers from the users on cli.
package cliquiz

import (
	"fmt"
	"reflect"

	"github.com/AlecAivazis/survey/v2"
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
