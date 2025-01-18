package list

import (
	"fmt"
	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"

	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
}

func NewCmdList(f factory.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available workspaces",
		Long: heredoc.Doc(`
				Retrieve and display a list of all workspaces associated 
				with your Asana account.`),
		Example: heredoc.Doc(`
				$ asana workspaces list
				$ asana workspaces ls
				$ asana ws ls
			`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runList(opts)
		},
	}

	return cmd
}

func runList(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

	client, err := opts.Client()
	if err != nil {
		return err
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	workspaces, err := client.AllWorkspaces()
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		fmt.Fprintf(opts.IO.Out, "No workspaces found for %s", cs.Bold(cfg.Username))
		return nil
	}

	fmt.Fprintf(opts.IO.Out, "\nWorkspaces of %s:\n\n", cs.Bold(cfg.Username))
	for i, ws := range workspaces {
		fmt.Fprintf(opts.IO.Out, "%d. %s\n", i+1, cs.Bold(ws.Name))
	}

	return nil
}
