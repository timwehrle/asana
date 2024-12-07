package login

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/timwehrle/act/api"
	"github.com/timwehrle/act/internal/auth"
	"github.com/timwehrle/act/internal/prompter"
	"github.com/timwehrle/act/internal/workspace"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your Asana account with your Personal Access Token.",
	Run: func(cmd *cobra.Command, args []string) {
		var token string

		token, err := prompter.Input("Please enter your Personal Access Token:", &token)

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

		client := api.New(token)

		workspaces, err := client.GetWorkspaces()
		if err != nil {
			fmt.Println("Error fetching workspaces:", err)
			return
		}

		if len(workspaces) == 0 {
			fmt.Println("No workspaces found.")
			return
		}

		options := make([]huh.Option[api.Workspace], len(workspaces))
		for i, workspace := range workspaces {
			options[i] = huh.Option[api.Workspace]{
				Key:   workspace.Name,
				Value: workspace,
			}
		}

		var selectedWorkspace api.Workspace

		huh.NewSelect[api.Workspace]().
			Title("Please select your default workspace:").
			Height(4).
			OptionsFunc(func() []huh.Option[api.Workspace] {
				return options
			}, nil).
			Value(&selectedWorkspace).
			Run()

		if err != nil {
			fmt.Println("Error selecting workspace:", err)
			return
		}

		err = workspace.SaveDefaultWorkspace(selectedWorkspace.GID)
		if err != nil {
			fmt.Println("Error saving default workspace:", err)
			return
		}

		fmt.Printf("Default workspace set to '%s'.\n", selectedWorkspace.Name)
	},
}
