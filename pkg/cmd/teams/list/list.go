package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type ListOptions struct {
	IO     *iostreams.IOStreams
	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
}

func NewCmdList(f factory.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		Client: f.Client,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all teams",
		Long: heredoc.Doc(`
				Retrieve and display a list of all teams assigned to your default workspace.
			`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF == nil {
				return listRun(opts)
			}
			return runF(opts)
		},
	}

	return cmd
}

func listRun(opts *ListOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	client, err := opts.Client()
	if err != nil {
		return err
	}

	teams, err := cfg.Workspace.AllTeams(client)
	if err != nil {
		return fmt.Errorf("failed to fetch teams: %w", err)
	}

	cs := opts.IO.ColorScheme()
	opts.IO.Printf("\nTeams in workspace %s:\n\n", cs.Bold(cfg.Workspace.Name))

	for i, team := range teams {
		opts.IO.Printf("%2d. %s\n", i+1, cs.Bold(team.Name))
	}

	return nil
}
