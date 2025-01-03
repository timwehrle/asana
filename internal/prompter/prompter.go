package prompter

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func Input(title, defaultValue string) (string, error) {
	var input string
	err := ask(&survey.Input{
		Message: title,
		Default: defaultValue,
	}, &input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func Confirm(message, defaultValue string) (bool, error) {
	var confirm bool
	err := ask(&survey.Confirm{
		Message: message,
		Default: defaultValue == "No",
	}, &confirm)
	return confirm, err
}

func Token() (string, error) {
	var token string
	err := ask(&survey.Password{
		Message: "Paste your authentication token:",
	}, &token, survey.WithValidator(survey.Required))
	return token, err
}

func Select(message string, options []string) (int, error) {
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

func Editor(message, existingDescription string) (string, error) {
	var input string

	err := ask(&survey.Editor{
		Message:       message,
		Default:       existingDescription,
		AppendDefault: true,
		HideDefault:   true,
	}, &input)
	if err != nil {
		return "", err
	}

	return input, nil
}

func ask(q survey.Prompt, response any, opts ...survey.AskOpt) error {
	opts = append(opts, survey.WithStdio(os.Stdin, os.Stdout, os.Stderr))
	if err := survey.AskOne(q, response, opts...); err != nil {
		return fmt.Errorf("could not prompt: %w", err)
	}
	return nil
}
