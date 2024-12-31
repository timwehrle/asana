package get

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/config/get/default_workspace"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a configuration value",
	}

	cmd.AddCommand(default_workspace.NewCmdDefaultWorkspace())

	return cmd
}
