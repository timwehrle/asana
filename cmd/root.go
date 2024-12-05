package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/timwehrle/act/cmd/auth"
)

var rootCmd = &cobra.Command{
	Use: "act",
}

func init() {
	rootCmd.AddCommand(auth.AuthCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
