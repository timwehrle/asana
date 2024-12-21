package auth

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/pkg/cmd/auth/login"
	"github.com/timwehrle/alfie/pkg/cmd/auth/logout"
	"github.com/timwehrle/alfie/pkg/cmd/auth/status"
)

func NewCmdAuth() *cobra.Command {
	cmd := &cobra.Command{
		Use: "auth <command>",
		Short: "Authenticate Alfie with Asana",
	}

	cmd.AddCommand(status.NewCmdStatus())
	cmd.AddCommand(login.NewCmdLogin())
	cmd.AddCommand(logout.NewCmdLogout())

	return cmd
}