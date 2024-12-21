package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/pkg/cmd/auth"
	"github.com/timwehrle/alfie/pkg/cmd/tasks"
)

var rootCmd = &cobra.Command{
	Use:   "alfie <command> <subcommand> [flags]",
	Short: "Alfie is a CLI tool for Asana",
	Long:  `Work with Asana from the command line.`,
}

func init() {
	rootCmd.AddCommand(auth.NewCmdAuth())
	rootCmd.AddCommand(tasks.NewCmdTasks())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
