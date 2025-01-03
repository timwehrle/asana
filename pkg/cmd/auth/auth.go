package auth

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/auth/login"
	"github.com/timwehrle/asana/pkg/cmd/auth/logout"
	"github.com/timwehrle/asana/pkg/cmd/auth/status"
	"github.com/timwehrle/asana/pkg/cmd/auth/update"
)

func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <subcommand>",
		Short: "Authenticate with Asana",
		Long: heredoc.Doc(`
			Manage authentication for the Asana CLI, including login
			logout and checking authentication status.
		`),
	}

	cmd.AddCommand(status.NewCmdStatus())
	cmd.AddCommand(login.NewCmdLogin())
	cmd.AddCommand(logout.NewCmdLogout())
	cmd.AddCommand(update.NewCmdUpdate())

	return cmd
}
