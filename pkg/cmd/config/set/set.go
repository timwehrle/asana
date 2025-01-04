package set

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/utils"
)

func NewCmdConfigSet(f factory.Factory) *cobra.Command {
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
			return runConfigSet(f, args[0])
		},
	}

	return cmd
}

func runConfigSet(f factory.Factory, key string) error {
	switch key {
	case "default-workspace", "dw":
		return setDefaultWorkspace(f)
	default:
		return fmt.Errorf("unknown configuration key: %s. Available keys are: default-workspace (dw)", key)
	}
}

func setDefaultWorkspace(f factory.Factory) error {
	client, err := f.NewAsanaClient()
	if err != nil {
		return err
	}

	workspaces, err := client.AllWorkspaces()
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

	cfg, err := f.Config()
	if err != nil {
		return err
	}

	// Workspace must be uppercase here since the Set function works with the
	// interface names and workspace is uppercased.
	err = cfg.Set("Workspace", selectedWorkspace)
	if err != nil {
		return err
	}

	fmt.Printf("%s Default workspace set to %s\n", utils.Success(), utils.Bold().Sprint(selectedWorkspace.Name))

	return nil
}
