package prompter

import (
	"github.com/charmbracelet/huh"
)

func Confirm(title string, defaultValue bool) (bool, error) {
	var confirm bool
	err := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&confirm).
		WithTheme(GlobalTheme).
		Run()
	return confirm, err
}
