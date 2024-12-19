package login

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/prompter"
	"github.com/timwehrle/alfie/internal/workspace"
)

var Cmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to an Asana account",
	Long: heredoc.Docf(`
			Authenticate with Asana using a Personal Access Token.
			
			To use this command, you'll need to generate a Personal Access Token from 
			your Asana account. Visit the settings in your Asana account, navigate to 
			the "Apps" tab, and create a new token. Ensure you store this token 
			securely and copy it for use in this command.

			Once you have the token, run the "login" command and provide the token 
			when prompted. This will save the token locally, enabling the application 
			to interact with the Asana API on your behalf. The token is securely stored 
			and used for authenticating subsequent API requests.

			If you encounter issues during the login process, double-check your token's 
			validity and ensure you have the necessary permissions granted for the 
			operations you intend to perform in Asana. If your token expires or is revoked, 
			you can generate a new one and repeat the login process.

			Note: Do not share your Personal Access Token with anyone, as it provides 
			full access to your account.
	`),
	Example: heredoc.Doc(`
		# Start login process
		$ act auth login
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		var token string

		_, err := auth.Get()
		if err == nil {
			fmt.Println("You are already logged in.")
			return nil
		}

		fmt.Print(heredoc.Doc(`
			Tip: You can generate a Personal Access Token here: https://app.asana.com/0/my-apps
		`))
		token, err = prompter.Token()
		if err != nil {
			return err
		}

		err = auth.Set(token)
		if err != nil {
			return err
		}

		fmt.Println("Successfully logged in.")

		client := api.New(token)

		workspaces, err := client.GetWorkspaces()
		if err != nil {
			return err
		}

		if len(workspaces) == 0 {
			fmt.Println("No workspaces found.")
			return nil
		}

		names := make([]string, len(workspaces))
		for i, ws := range workspaces {
			names[i] = ws.Name
		}

		index, err := prompter.Select("Please select your default workspace:", names)
		if err != nil {
			return err
		}

		selectedWorkspace := workspaces[index]

		err = workspace.SaveDefaultWorkspace(selectedWorkspace.GID, selectedWorkspace.Name)
		if err != nil {
			return err
		}

		fmt.Printf("Default workspace set to '%s'.\n", selectedWorkspace.Name)

		return nil
	},
}
