package status

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/spf13/cobra"
)

type Status struct {
	LoggedIn       bool
	APIOperational bool
	User           *asana.User
	WorkspaceID    string
	WorkspaceName  string
}

type StatusOptions struct {
	factory.Factory
	IO *iostreams.IOStreams
}

func NewCmdStatus(f factory.Factory) *cobra.Command {
	opts := &StatusOptions{
		Factory: f,
		IO:      f.IOStreams(),
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "View current authentication status",
		Long: heredoc.Doc(`
            Display the current authentication status, including:
            - Login state
            - API connectivity
            - User information
            - Default workspace configuration
            
            This command helps verify your setup and connectivity to Asana.`),
		Example: heredoc.Docf(`
            # Check authentication status
            $ %[1]s auth status
            
            # View status with debug information
            $ %[1]s auth status --debug
        `, "asana"),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(opts)
		},
	}

	return cmd
}

func runStatus(opts *StatusOptions) error {
	status, err := getStatus(opts)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	return printStatus(opts.IO, status)
}

func getStatus(opts *StatusOptions) (*Status, error) {
	status := &Status{}

	token, err := auth.Get()
	if err != nil {
		status.LoggedIn = false
		return status, nil
	}
	status.LoggedIn = token != ""

	cfg, err := opts.Factory.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	status.WorkspaceID = cfg.Workspace.ID
	status.WorkspaceName = cfg.Workspace.Name

	if status.LoggedIn {
		client, err := opts.NewAsanaClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create Asana client: %w", err)
		}

		user, err := client.CurrentUser()
		if err != nil {
			status.APIOperational = false
			return status, nil
		}

		status.APIOperational = true
		status.User = user
	}

	return status, nil
}

func printStatus(io *iostreams.IOStreams, status *Status) error {
	cs := io.ColorScheme()

	if !status.LoggedIn {
		fmt.Fprintf(io.Out, "%s %s\n", cs.WarningIcon, cs.Bold("Not logged in"))
		return nil
	}

	if status.APIOperational {
		fmt.Fprintf(io.Out, "%s %s\n", cs.SuccessIcon, cs.Bold("API is operational"))
	} else {
		fmt.Fprintf(io.Out, "%s %s\n", cs.ErrorIcon, cs.Bold("API is not responding"))
	}

	if status.User != nil {
		fmt.Fprintf(io.Out, "\n%s:\n", cs.Bold("User Information"))
		fmt.Fprintf(io.Out, "  Name: %s\n", status.User.Name)
		fmt.Fprintf(io.Out, "  ID:   %s\n", status.User.ID)
	}

	fmt.Fprintf(io.Out, "\n%s:\n", cs.Bold("Workspace Configuration"))
	if status.WorkspaceID == "" || status.WorkspaceName == "" {
		fmt.Fprintf(io.Out, "%s %s\n", cs.WarningIcon, cs.Bold("No default workspace configured"))
	} else {
		fmt.Fprintf(io.Out, "  Name: %s\n", status.WorkspaceName)
		fmt.Fprintf(io.Out, "  ID:   %s\n", status.WorkspaceID)
	}

	return nil
}
