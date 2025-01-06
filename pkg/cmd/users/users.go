package users

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/users/list"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdUsers(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users of your Asana workspace",
		Long:  "Manage users of your Asana workspace",
	}

	cmd.AddCommand(list.NewCmdList(f))

	return cmd
}
