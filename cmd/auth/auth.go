package auth

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/jodot/cmd/auth/login"
	"github.com/timwehrle/jodot/cmd/auth/logout"
	"github.com/timwehrle/jodot/cmd/auth/status"
)

var AuthCmd = &cobra.Command{
	Use:   "auth <command>",
	Short: "Authenticate jodot with Asana",
}

func init() {
	AuthCmd.AddCommand(login.LoginCmd)
	AuthCmd.AddCommand(logout.LogoutCmd)
	AuthCmd.AddCommand(status.StatusCmd)
}
