package tasks

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/tasks/list"
	"github.com/timwehrle/asana/pkg/cmd/tasks/update"
	"github.com/timwehrle/asana/pkg/cmd/tasks/view"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdTasks(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tasks <subcommand>",
		Aliases: []string{"ts"},
		Short:   "Manage your Asana tasks",
		Long:    "Perform operations related to your Asana tasks.",
	}

	cmd.AddCommand(list.NewCmdList(f))
	cmd.AddCommand(view.NewCmdView(f))
	cmd.AddCommand(update.NewCmdUpdate(f))

	return cmd
}
