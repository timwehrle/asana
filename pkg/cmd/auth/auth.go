package auth

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/auth/login"
	"github.com/timwehrle/asana/pkg/cmd/auth/logout"
	"github.com/timwehrle/asana/pkg/cmd/auth/status"
	"github.com/timwehrle/asana/pkg/cmd/auth/update"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdAuth(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <subcommand>",
		Short: "Authenticate with Asana",
		Long: heredoc.Doc(`
			Manage authentication for the Asana CLI, including login
			logout and checking authentication status.`),
	}

	cmd.AddCommand(status.NewCmdStatus(f, nil))
	cmd.AddCommand(login.NewCmdLogin(f, nil))
	cmd.AddCommand(logout.NewCmdLogout(f, nil))
	cmd.AddCommand(update.NewCmdUpdate(f, nil))

	return cmd
}
