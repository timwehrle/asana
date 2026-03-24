package comments

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/tasks/comments/add"
	"github.com/timwehrle/asana/pkg/cmd/tasks/comments/list"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdComments(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comments <subcommand>",
		Short: "Read and add comments on a task",
		Long:  "List comment stories on a task or add a new comment.",
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(add.NewCmdAdd(f, nil))

	return cmd
}
