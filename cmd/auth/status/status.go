package status

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/workspace"
)

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of the user.",
	Long: heredoc.Doc(`
			Get the status of the current user and API.

			This command displays the API's operational status, 
			the logged-in user's username, and the default workspace.
	`),
	Example: heredoc.Doc(`
			# Display status 
			$ act auth status
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		gid, name, err := workspace.LoadDefaultWorkspace()
		if err != nil {
			return err
		}

		token, err := auth.Get()
		if err != nil {
			fmt.Println("You are not logged in.")
			return nil
		}

		client := api.New(token)

		me, err := client.GetMe()
		if err != nil {
			return err
		} else {
			fmt.Println("API is operational.")
		}

		fmt.Printf("Logged in as: %s (%s)\n", me.Name, me.GID)
		if gid == "" || name == "" {
			fmt.Println("No default workspace set.")
		} else {
			fmt.Printf("Default workspace: %s (%s)\n", name, gid)
		}

		return nil
	},
}
