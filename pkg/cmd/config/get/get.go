package get

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/utils"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print the value of a given configuration key",
		Example: heredoc.Doc(`
				# Get a configuration value
				$ asana config get default-workspace
				$ asana config get dw
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigGet(args[0])
		},
	}

	return cmd
}

func runConfigGet(key string) error {
	switch key {
	case "default-workspace", "dw":
		defaultWorkspace, err := config.GetDefaultWorkspace()
		if err != nil {
			return err
		}

		fmt.Printf("Default workspace is %s (%s)\n", utils.Bold().Sprint(defaultWorkspace.Name), defaultWorkspace.ID)
		return nil

	default:
		return fmt.Errorf("unknown configuration key: %s. Available keys are: default-workspace (dw)", key)
	}
}
