package login

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/timwehrle/act/internal/auth"
	"github.com/timwehrle/act/internal/prompter"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your Asana account with your Personal Access Token.",
	Run: func(cmd *cobra.Command, args []string) {
		var token string

		err := huh.NewInput().
			Title("Enter your Personal Access Token:").
			Value(&token).
			WithTheme(prompter.GlobalTheme).
			Run()

		if err != nil {
			fmt.Println("Error reading token:", err)
			return
		}

		//! Method won't work in WSL2/Linux since it is a bug within zalando/go-keyring
		err = auth.Set(token)
		if err != nil {
			fmt.Println("Error storing credentials:", err)
			return
		}

		fmt.Println("Successfully logged in.")
	},
}
