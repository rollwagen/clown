package prompter

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/rollwagen/clown/pkg/prompter"
)

func New() prompter.Prompter {
	return &surveyPrompter{}
}

type surveyPrompter struct{}

func (p *surveyPrompter) Select(message, defaultValue string, options []string) (result int, err error) {
	q := &survey.Select{
		Message:  message,
		Options:  options,
		PageSize: 15,
	}

	if defaultValue != "" {
		q.Default = defaultValue
	}

	err = survey.AskOne(q, &result)
	if err != nil {
		return -1, err
	}

	return result, nil
}
