package status

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/act/internal/workspace"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of the user.",
	Run: func(cmd *cobra.Command, args []string) {
		workspace, err := workspace.LoadDefaultWorkspace()
		if err != nil {
			fmt.Println("Error fetching default workspace:", err)
			return
		}

		fmt.Println("Your default workspace is set to:", workspace)
	},
}
