package tags

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/tags/list"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdTags(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage your Asana tags",
		Long:  "Perform operations related to your Asana tags.",
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
