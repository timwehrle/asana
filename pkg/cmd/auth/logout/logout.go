package logout

import (
	"fmt"
	"github.com/timwehrle/asana/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/prompter"
)

func NewCmdLogout() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of your Asana account",
		Long: heredoc.Doc(`
				Log out of your current Asana account by removing locally 
				stored credentials.
				
				This action revokes CLI access to the Asana API.`),
		Example: heredoc.Doc(`$ asana auth logout`),
		RunE: func(_ *cobra.Command, _ []string) error {
			return logoutRun()
		},
	}

	return cmd
}

func logoutRun() error {
	_, err := auth.Get()
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
		fmt.Println(utils.Success(), "Logged out")
	} else {
		fmt.Println("Logout aborted.")
	}

	return nil
}
