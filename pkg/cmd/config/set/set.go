package set

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/config/set/default_workspace"
)

func NewCmdSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a configuration value",
	}

	cmd.AddCommand(default_workspace.NewCmdDefaultWorkspace())

	return cmd
}
