package prompter

import "github.com/charmbracelet/huh"

func Input(title string, value *string) (string, error) {
	var input string
	err := huh.NewInput().
		Title(title).
		Value(&input).
		WithTheme(GlobalTheme).
		Run()
	return input, err
}
