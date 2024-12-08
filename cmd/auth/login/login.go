package login

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
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

		fmt.Print(heredoc.Doc(`
			Tip: You can generate a Personal Access Token here: https://app.asana.com/0/my-apps
		`))
		token, err := prompter.Token()
		if err != nil {
			fmt.Println("Error fetching token:", err)
			return
		}

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

		names := make([]string, len(workspaces))
		for i, ws := range workspaces {
			names[i] = ws.Name
		}

		prompt := &survey.Select{
			Message: "Please select your default workspace:",
			Options: names,
		}

		answerIndex := 0

		err = survey.AskOne(prompt, &answerIndex)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		selectedWorkspace := workspaces[answerIndex]

		err = workspace.SaveDefaultWorkspace(selectedWorkspace.GID, selectedWorkspace.Name)
		if err != nil {
			fmt.Println("Error saving default workspace:", err)
			return
		}

		fmt.Printf("Default workspace set to '%s'.\n", selectedWorkspace.Name)
	},
}
