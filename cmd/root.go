package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/timwehrle/alaric/cmd/auth"
	"github.com/timwehrle/alaric/cmd/tasks"
)

var rootCmd = &cobra.Command{
	Use: "alaric",
}

func init() {
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(tasks.TasksCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
