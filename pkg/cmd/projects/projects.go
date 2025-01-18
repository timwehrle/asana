package projects

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/projects/list"
	"github.com/timwehrle/asana/pkg/cmd/projects/tasks"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdProjects(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects <subcommand>",
		Short: "Manage your Asana projects",
		Long:  "Perform operations related to your Asana projects.",
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(tasks.NewCmdTasks(f, nil))

	return cmd
}
