package tasks

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/tasks/list"
	"github.com/timwehrle/asana/pkg/cmd/tasks/view"
)

func NewCmdTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Manage your Asana tasks",
		Long:  "Perform operations related to your Asana tasks.",
	}

	cmd.AddCommand(list.NewCmdList())
	cmd.AddCommand(view.NewCmdView())

	return cmd
}
