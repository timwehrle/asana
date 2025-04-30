package teams

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/teams/list"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdTeams(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "teams <subcommand>",
		Short: "Manage your Asana teams",
		Long:  "Perform operations related to your Asana teams.",
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
