package auth

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/cmd/auth/login"
	"github.com/timwehrle/alfie/cmd/auth/logout"
	"github.com/timwehrle/alfie/cmd/auth/status"
)

var Cmd = &cobra.Command{
	Use:   "auth <command>",
	Short: "Authenticate alfie with Asana",
}

func init() {
	Cmd.AddCommand(login.Cmd)
	Cmd.AddCommand(logout.Cmd)
	Cmd.AddCommand(status.Cmd)
}
