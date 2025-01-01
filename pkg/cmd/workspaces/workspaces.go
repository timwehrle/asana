package workspaces

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/workspaces/list"
)

func NewCmdWorkspace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspaces",
		Short: "Manage your Asana workspaces",
		Long: heredoc.Doc(`
				Perform operations related to your Asana workspaces.
		`),
	}

	cmd.AddCommand(list.NewCmdList())

	return cmd
}
