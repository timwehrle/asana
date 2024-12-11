package auth

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/alaric/cmd/auth/login"
	"github.com/timwehrle/alaric/cmd/auth/logout"
	"github.com/timwehrle/alaric/cmd/auth/status"
)

var AuthCmd = &cobra.Command{
	Use:   "auth <command>",
	Short: "Authenticate ACT with Asana",
}

func init() {
	AuthCmd.AddCommand(login.LoginCmd)
	AuthCmd.AddCommand(logout.LogoutCmd)
	AuthCmd.AddCommand(status.StatusCmd)
}
