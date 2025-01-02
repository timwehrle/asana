package status

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/utils"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
)

func NewCmdStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "View current authentication status",
		Long: heredoc.Doc(`
				Display the current authentication status, including 
				the logged-in user and API health. This command helps 
				verify connectivity and user identity.
		`),
		Example: heredoc.Doc(`
				# Check authentication status
				$ asana auth status
		`),
		RunE: func(_ *cobra.Command, _ []string) error {
			return statusRun()
		},
	}

	return cmd
}

func statusRun() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	token, err := auth.Get()
	if err != nil {
		fmt.Println("You are not logged in.")
		return nil
	}

	client := asana.NewClientWithAccessToken(token)

	me, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("API is operational.")
	fmt.Printf("Logged in as %s (%s)\n", utils.Bold().Sprintf(me.Name), me.ID)

	if cfg.Workspace.ID == "" || cfg.Workspace.Name == "" {
		fmt.Println("No default workspace set.")
	} else {
		fmt.Printf("Default workspace is %s (%s)\n", utils.Bold().Sprintf(cfg.Workspace.Name), cfg.Workspace.ID)
	}

	return nil
}
