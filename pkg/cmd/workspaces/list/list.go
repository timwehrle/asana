package list

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/utils"
)

func NewCmdList(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available workspaces",
		Long: heredoc.Doc(`
				Retrieve and display a list of all workspaces associated 
				with your Asana account.
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

	workspaces, err := client.AllWorkspaces()
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
		return nil
	}

	fmt.Println(utils.BoldUnderline().Sprint("Your Workspaces:"))
	for i, ws := range workspaces {
		fmt.Printf("%d. %s\n", i+1, ws.Name)
	}

	return nil
}
