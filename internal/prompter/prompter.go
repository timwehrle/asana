package prompter

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type Prompter interface {
	Input(prompt, defaultValue string) (string, error)
	Confirm(prompt, defaultValue string) (bool, error)
	Token() (string, error)
	Select(message string, options []string) (int, error)
	Editor(prompt, existingDescription string) (string, error)
}

type DefaultPrompter struct{}

func New() *DefaultPrompter {
	return &DefaultPrompter{}
}

func (p *DefaultPrompter) Input(prompt, defaultValue string) (string, error) {
	var result string
	err := ask(&survey.Input{
		Message: prompt,
		Default: defaultValue,
	}, &result)
	if err != nil {
		return "", err
	}

	result = strings.TrimSpace(result)

	if result == "" {
		return defaultValue, err
	}

	return result, nil
}

func (p *DefaultPrompter) Confirm(prompt, defaultValue string) (bool, error) {
	var confirm bool
	err := ask(&survey.Confirm{
		Message: prompt,
		Default: defaultValue == "No",
	}, &confirm)
	return confirm, err
}

func (p *DefaultPrompter) Token() (string, error) {
	var token string
	err := ask(&survey.Password{
		Message: "Paste your authentication token:",
	}, &token, survey.WithValidator(survey.Required))
	return token, err
}

func (p *DefaultPrompter) Select(message string, options []string) (int, error) {
	var answerIndex int

	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	err := ask(prompt, &answerIndex)
	if err != nil {
		return -1, err
	}

	return answerIndex, nil
}

func (p *DefaultPrompter) Editor(prompt, existingDescription string) (string, error) {
	var input string

	err := ask(&survey.Editor{
		Message:       prompt,
		Default:       existingDescription,
		AppendDefault: true,
		HideDefault:   true,
	}, &input)
	if err != nil {
		return "", err
	}

	return input, nil
}

var ask = func(q survey.Prompt, response any, opts ...survey.AskOpt) error {
	opts = append(opts, survey.WithStdio(os.Stdin, os.Stdout, os.Stderr))
	err := survey.AskOne(q, response, opts...)
	if err == nil {
		return nil
	}
	return fmt.Errorf("could not prompt: %w", err)
}
