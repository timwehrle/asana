package set

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/api"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/utils"
)

func NewCmdConfigSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key>",
		Short: "Update configuration with a value",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
				# Set a configuration value
				$ asana config set default-workspace
				$ asana config set dw
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSet(args[0])
		},
	}

	return cmd
}

func runConfigSet(key string) error {
	switch key {
	case "default-workspace", "dw":
		return setDefaultWorkspace()
	default:
		return fmt.Errorf("unknown configuration key: %s. Available keys are: default-workspace (dw)", key)
	}
}

func setDefaultWorkspace() error {
	token, err := auth.Get()
	if err != nil {
		return err
	}

	client := api.New(token)

	workspaces, err := client.GetWorkspaces()
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
		return nil
	}

	names := make([]string, len(workspaces))
	for i, ws := range workspaces {
		names[i] = ws.Name
	}

	index, err := prompter.Select("Select a new default workspace:", names)
	if err != nil {
		return err
	}

	selectedWorkspace := workspaces[index]

	err = config.UpdateDefaultWorkspace(selectedWorkspace.GID, selectedWorkspace.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s Default workspace set to %s\n", utils.Success(), utils.Bold().Sprintf(selectedWorkspace.Name))

	return nil
}
