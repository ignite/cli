// Package cliquiz is a tool to collect answers from the users on cli.
package cliquiz

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

type Question struct {
	Question      string
	DefaultAnswer interface{}
	Answer        interface{}
}

func NewQuestion(question string, defaultAnswer, answer interface{}) Question {
	return Question{question, defaultAnswer, answer}
}

func Ask(question ...Question) error {
	for _, q := range question {
		prompt := promptui.Prompt{
			Label:   q.Question,
			Default: fmt.Sprintf("%v", q.DefaultAnswer),
		}

		result, err := prompt.Run()
		if err != nil {
			return err
		}
		if _, ok := q.Answer.(*string); ok {
			result = strconv.Quote(result)
		}
		// instead of going with the type switches to convert result to propver type,
		// we can let this big work to json package.
		if err := json.Unmarshal([]byte(result), q.Answer); err != nil {
			return err
		}
	}
	return nil
}
