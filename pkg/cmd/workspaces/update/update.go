package update

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/api"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/utils"
)

func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateRun()
		},
	}

	return cmd
}

func updateRun() error {
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

	index, err := prompter.Select("Please select the the new default workspace:", names)
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
