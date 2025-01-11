package list

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	factory.Factory
	IO *iostreams.IOStreams
}

func NewCmdList(f factory.Factory) *cobra.Command {
	opts := &ListOptions{
		Factory: f,
		IO:      f.IOStreams(),
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
			return listRun(opts)
		},
	}

	return cmd
}

func listRun(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

	client, err := opts.NewAsanaClient()
	if err != nil {
		return err
	}

	cfg, err := opts.Factory.Config()
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
