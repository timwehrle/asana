package status

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/api"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
)

func NewCmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display the current user's status and API health",
		Long: heredoc.Doc(`
				Display the status of the current logged-in user and the API.

				This command shows whether the API is running, the current
				user and the default config.
		`),
		Example: heredoc.Doc(`
				# Start status process
				$ asana auth status
		`),
		RunE: func(_ *cobra.Command, _ []string) error {
			return statusRun()
		},
	}

	return cmd
}

func statusRun() error {
	config, err := config.LoadConfig()
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
	}
	fmt.Println("API is operational.")

	fmt.Printf("Logged in as: %s (%s)\n", me.Name, me.GID)
	if config.Workspace.GID == "" || config.Workspace.Name == "" {
		fmt.Println("No default config set.")
	} else {
		fmt.Printf("Default config: %s (%s)\n", config.Workspace.Name, config.Workspace.GID)
	}

	return nil
}
