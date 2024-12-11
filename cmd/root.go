package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/timwehrle/jodot/cmd/auth"
	"github.com/timwehrle/jodot/cmd/tasks"
)

var rootCmd = &cobra.Command{
	Use: "jodot",
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
