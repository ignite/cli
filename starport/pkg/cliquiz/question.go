// Package cliquiz is a tool to collect answers from the users on cli.
package cliquiz

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type Question struct {
	Question      string
	DefaultAnswer interface{}
	Answer        interface{}
}

func NewQuestion(question string, defaultAnswer, answer interface{}) Question {
	return Question{question, defaultAnswer, answer}
}

// Ask asks questions and collect answers.
func Ask(question ...Question) error {
	for _, q := range question {
		prompt := &survey.Input{
			Message: q.Question,
			Default: fmt.Sprintf("%v", q.DefaultAnswer),
		}

		if err := survey.AskOne(prompt, q.Answer); err != nil {
			return err
		}
	}
	return nil
}
