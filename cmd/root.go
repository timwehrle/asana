package cmd

import (
	"github.com/timwehrle/asana/internal/version"
	"github.com/timwehrle/asana/pkg/factory"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	service "github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/pkg/cmd/auth"
	"github.com/timwehrle/asana/pkg/cmd/config"
	"github.com/timwehrle/asana/pkg/cmd/tasks"
	"github.com/timwehrle/asana/pkg/cmd/workspaces"
)

var rootCmd = &cobra.Command{
	Use:     "asana <command> <subcommand> [flags]",
	Short:   "The Asana CLI tool",
	Version: version.Version,
	Long:    `Work with Asana from the command line.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.HasParent() && cmd.Parent().Name() == "auth" {
			return nil
		}

		err := service.Check()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cmdFactory := factory.New()

	rootCmd.AddCommand(auth.NewCmdAuth(cmdFactory))
	rootCmd.AddCommand(tasks.NewCmdTasks(cmdFactory))
	rootCmd.AddCommand(workspaces.NewCmdWorkspace(cmdFactory))
	rootCmd.AddCommand(config.NewCmdConfig(cmdFactory))

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	// Colorize output
	rootCmd.SetOut(color.Output)
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgGreen).SprintFunc())
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Examples:`, `{{StyleHeading "Examples:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Flags:`, `{{StyleHeading "Flags:"}}`,
	).Replace(usageTemplate)
	re := regexp.MustCompile(`(?m)^Flags:\s*$`)
	usageTemplate = re.ReplaceAllLiteralString(usageTemplate, `{{StyleHeading "Flags:"}}`)
	rootCmd.SetUsageTemplate(usageTemplate)
}

func Execute() error {
	return rootCmd.Execute()
}
