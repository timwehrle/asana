package logout

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/act/internal/auth"
	"github.com/timwehrle/act/internal/prompter"
)

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your Asana account.",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.Get()
		if err == auth.ErrNoToken {
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
