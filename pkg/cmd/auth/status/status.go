package status

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/utils"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
)

func NewCmdStatus(f factory.Factory) *cobra.Command {
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
			return statusRun(f)
		},
	}

	return cmd
}

func statusRun(f factory.Factory) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	_, err = auth.Get()
	if err != nil {
		fmt.Println("You are not logged in.")
		return nil
	}

	client, err := f.NewAsanaClient()
	if err != nil {
		return err
	}

	me, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("API is operational.")
	fmt.Printf("Logged in as %s (%s)\n", utils.Bold().Sprint(me.Name), me.ID)

	if cfg.Workspace.ID == "" || cfg.Workspace.Name == "" {
		fmt.Println("No default workspace set.")
	} else {
		fmt.Printf("Default workspace is %s (%s)\n", utils.Bold().Sprint(cfg.Workspace.Name), cfg.Workspace.ID)
	}

	return nil
}
