package list

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/utils"

	"github.com/spf13/cobra"
)

func NewCmdList(f factory.Factory) *cobra.Command {
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
			return listRun(f)
		},
	}

	return cmd
}

func listRun(f factory.Factory) error {
	client, err := f.NewAsanaClient()
	if err != nil {
		return err
	}

	cfg, err := f.Config()
	if err != nil {
		return err
	}

	workspaces, err := client.AllWorkspaces()
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		fmt.Printf("No workspaces found for %s", utils.Bold().Sprint(cfg.Username))
		return nil
	}

	fmt.Printf("\nWorkspaces of %s:\n\n", utils.Bold().Sprint(cfg.Username))
	for i, ws := range workspaces {
		fmt.Printf("%d. %s\n", i+1, utils.Bold().Sprint(ws.Name))
	}

	return nil
}
