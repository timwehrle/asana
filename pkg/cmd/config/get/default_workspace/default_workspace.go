package default_workspace

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/utils"
)

func NewCmdDefaultWorkspace() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "default-workspace",
		Aliases: []string{"dw"},
		Short:   "Get the default workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDefaultWorkspace()
		},
	}

	return cmd
}

func runDefaultWorkspace() error {
	defaultWorkspace, err := config.GetDefaultWorkspace()
	if err != nil {
		return err
	}

	fmt.Printf("Default workspace is %s (%s)\n", utils.Bold().Sprintf(defaultWorkspace.Name), defaultWorkspace.GID)

	return nil
}
