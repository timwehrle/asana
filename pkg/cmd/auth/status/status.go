package status

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/workspace"
)

func NewCmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use: "status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return statusRun()
		},
	}

	return cmd
}

func statusRun() error {
			gid, name, err := workspace.LoadDefaultWorkspace()
		if err != nil {
			return err
		}

		token, err := auth.Get()
		if err != nil {
			fmt.Println("You are not logged in.")
			return nil
		}

		client := api.New(token)

		me, err := client.GetMe()
		if err != nil {
			return err
		} else {
			fmt.Println("API is operational.")
		}

		fmt.Printf("Logged in as: %s (%s)\n", me.Name, me.GID)
		if gid == "" || name == "" {
			fmt.Println("No default workspace set.")
		} else {
			fmt.Printf("Default workspace: %s (%s)\n", name, gid)
		}

		return nil
}
