package logout

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/prompter"
)

func NewCmdLogout() *cobra.Command {
	cmd := &cobra.Command{
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
			$ asana auth logout
		`),
		RunE: func(_ *cobra.Command, _ []string) error {
			return logoutRun()
		},
	}

	return cmd
}

func logoutRun() error {
	_, err := auth.Get()
	if errors.Is(err, auth.ErrNoToken) {
		fmt.Println("No active session found. You are not logged in.")
		return nil
	}
	if err != nil {
		return err
	}

	confirm := false
	confirm, err = prompter.Confirm("Are you sure you want to log out?", "No")
	if err != nil {
		return err
	}

	if confirm {
		err := auth.Delete()
		if err != nil {
			return err
		}
		fmt.Println("Successfully logged out.")
	} else {
		fmt.Println("Logout aborted.")
	}

	return nil
}
