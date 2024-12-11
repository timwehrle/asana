package auth

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/cmd/auth/login"
	"github.com/timwehrle/alfie/cmd/auth/logout"
	"github.com/timwehrle/alfie/cmd/auth/status"
)

var AuthCmd = &cobra.Command{
	Use:   "auth <command>",
	Short: "Authenticate alfie with Asana",
}

func init() {
	AuthCmd.AddCommand(login.LoginCmd)
	AuthCmd.AddCommand(logout.LogoutCmd)
	AuthCmd.AddCommand(status.StatusCmd)
}
