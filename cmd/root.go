package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/pkg/cmd/auth"
	"github.com/timwehrle/alfie/pkg/cmd/brief"
	"regexp"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "alfie <command> <subcommand> [flags]",
	Short: "Alfie is a CLI tool for Asana",
	Long:  `Work with Asana from the command line.`,
}

func init() {
	rootCmd.AddCommand(auth.NewCmdAuth())
	rootCmd.AddCommand(brief.NewCmdBrief())

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
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}
