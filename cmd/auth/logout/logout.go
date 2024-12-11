package logout

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/prompter"
)

var Cmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of an Asana account",
	Long: heredoc.Doc(`
			Log out of the current Asana account.

			This command will remove the locally stored credentials,
			disabling the application from interacting with the Asana API.

			This command does not invalidate the Personal Access Token.

			Note: This action is irreversible. If you log out, you will need to
			repeat the login process to regain access to the Asana API.
	`),
	Example: heredoc.Doc(`
		# Start logout process
		$ act auth logout
	`),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.Get()
		if errors.Is(err, auth.ErrNoToken) {
			fmt.Println("No active session found. You are not logged in.")
			return
		}
		if err != nil {
			fmt.Println("Error retrieving user details:", err)
			return
		}

		confirm := false
		confirm, err = prompter.Confirm("Are you sure you want to log out?", "No")
		if err != nil {
			fmt.Println("Error reading confirmation:", err)
			return
		}

		if confirm {
			err := auth.Delete()
			if err != nil {
				fmt.Println("Error deleting credentials:", err)
				return
			}

			fmt.Println("Successfully logged out.")
		} else {
			fmt.Println("Logout aborted.")
		}
	},
}
