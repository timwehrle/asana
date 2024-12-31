package workspaces

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/workspaces/list"
	"github.com/timwehrle/asana/pkg/cmd/workspaces/update"
)

func NewCmdWorkspace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspaces",
		Short: "Work with your workspaces",
	}

	cmd.AddCommand(list.NewCmdList())
	cmd.AddCommand(update.NewCmdUpdate())

	return cmd
}
