package auth

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/auth/login"
	"github.com/timwehrle/asana/pkg/cmd/auth/logout"
	"github.com/timwehrle/asana/pkg/cmd/auth/status"
)

func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate the CLI with Asana",
	}

	cmd.AddCommand(status.NewCmdStatus())
	cmd.AddCommand(login.NewCmdLogin())
	cmd.AddCommand(logout.NewCmdLogout())

	return cmd
}
