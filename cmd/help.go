package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

// HelpSection represents a named section in the help output
type HelpSection struct {
	Title string
	Body  string
}

// suggestNestedCommands provides suggestions for unknown commands
func suggestNestedCommands(w io.Writer, cmd *cobra.Command, arg string) {
	fmt.Fprintf(w, "unknown command %q for %q\n", arg, cmd.CommandPath())

	suggestions := getSuggestions(cmd, arg)
	if len(suggestions) > 0 {
		fmt.Fprint(w, "\nDid you mean this?\n")
		for _, suggestion := range suggestions {
			fmt.Fprintf(w, "\t%s\n", suggestion)
		}
	}

	fmt.Fprint(w, "\n")
	_ = showRootUsage(cmd)
}

// getSuggestions returns command suggestions based on input
func getSuggestions(cmd *cobra.Command, arg string) []string {
	if arg == "help" {
		return []string{"--help"}
	}

	if cmd.SuggestionsMinimumDistance <= 0 {
		cmd.SuggestionsMinimumDistance = 2
	}
	return cmd.SuggestionsFor(arg)
}

// showHelp displays the help content for a command
func showHelp(cmd *cobra.Command, _ []string, w io.Writer) {
	flags := cmd.Flags()

	if help, _ := flags.GetBool("help"); !help && !cmd.Runnable() && len(flags.Args()) > 0 {
		suggestNestedCommands(os.Stderr, cmd, flags.Args()[0])
		return
	}

	sections := buildHelpSections(cmd)
	for _, section := range sections {
		printSection(w, section)
	}
}

// printSection formats and prints a help section
func printSection(w io.Writer, section HelpSection) {
	io := iostreams.System()
	cs := io.ColorScheme()

	if section.Title != "" {
		fmt.Fprintln(w, cs.Bold(section.Title))
		fmt.Fprintln(w, format.Indent(strings.Trim(section.Body, "\r\n"), "  "))
	} else {
		fmt.Fprintln(w, section.Body)
	}
	fmt.Fprintln(w)
}

// buildHelpSections creates all help sections for a command
func buildHelpSections(cmd *cobra.Command) []HelpSection {
	var sections []HelpSection

	// Add description section
	if desc := getCommandDescription(cmd); desc != "" {
		sections = append(sections, HelpSection{"", desc})
	}

	// Add usage section
	sections = append(sections, HelpSection{"Usage", cmd.UseLine()})

	// Add aliases section if present
	if len(cmd.Aliases) > 0 {
		sections = append(sections, HelpSection{"Aliases", formatAliases(cmd)})
	}

	// Add available commands if any
	if commands := cmd.Commands(); len(commands) > 0 {
		sections = append(sections, HelpSection{"Available Commands", formatCommands(commands)})
	}

	// Add flags sections
	appendFlagSections(cmd, &sections)

	// Add examples if present
	if cmd.Example != "" {
		sections = append(sections, HelpSection{"Examples", cmd.Example})
	}

	// Add learn more section
	sections = append(sections, getLearnMoreSection())

	return sections
}

// getCommandDescription returns the long or short description
func getCommandDescription(cmd *cobra.Command) string {
	if cmd.Long != "" {
		return cmd.Long
	}
	return cmd.Short
}

// appendFlagSections adds local and inherited flags sections
func appendFlagSections(cmd *cobra.Command, sections *[]HelpSection) {
	if cmd.HasAvailableLocalFlags() {
		*sections = append(*sections, HelpSection{"Flags", format.Dedent(cmd.LocalFlags().FlagUsages())})
	}

	if inheritedFlags := cmd.InheritedFlags().FlagUsages(); inheritedFlags != "" {
		*sections = append(*sections, HelpSection{"Inherited Flags", format.Dedent(inheritedFlags)})
	}
}

// formatAliases formats command aliases with parent aliases
func formatAliases(cmd *cobra.Command) string {
	return strings.Join(getParentAliases(cmd, cmd.Aliases), ", ")
}

// getParentAliases recursively builds alias list including parent commands
func getParentAliases(cmd *cobra.Command, aliases []string) []string {
	if !cmd.HasParent() {
		return aliases
	}

	parentAliases := append(cmd.Parent().Aliases, cmd.Parent().Name())
	sort.Strings(parentAliases)

	combinedAliases := make([]string, 0, len(aliases)*len(parentAliases))
	for _, alias := range aliases {
		for _, parentAlias := range parentAliases {
			combinedAliases = append(combinedAliases, fmt.Sprintf("%s %s", parentAlias, alias))
		}
	}

	return getParentAliases(cmd.Parent(), combinedAliases)
}

// formatCommands formats the list of available commands
func formatCommands(commands []*cobra.Command) string {
	maxLength := getMaxCommandLength(commands)
	var sb strings.Builder

	for _, cmd := range commands {
		sb.WriteString(fmt.Sprintf("%s %s\n", padRight(cmd.Name(), maxLength), cmd.Short))
	}

	return sb.String()
}

// getMaxCommandLength calculates the maximum command name length
func getMaxCommandLength(commands []*cobra.Command) int {
	maxLen := 0
	for _, cmd := range commands {
		if length := len(cmd.Name()); length > maxLen {
			maxLen = length
		}
	}
	return maxLen + 2
}

// padRight pads a string with spaces to the specified length
func padRight(s string, length int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", length), s)
}

// getLearnMoreSection returns the help section with additional information
func getLearnMoreSection() HelpSection {
	return HelpSection{
		Title: "Learn More",
		Body: heredoc.Docf(`
            Use %[1]sasana <command> <subcommand> --help%[1]s for more information about a command.
        `, "`"),
	}
}

// showRootUsage displays the root command usage
func showRootUsage(cmd *cobra.Command) error {
	io := iostreams.System()
	cs := io.ColorScheme()

	fmt.Print(cs.Bold("Usage:"))
	fmt.Printf("  %s <command> <subcommand> [flags]\n", cmd.Name())

	fmt.Println("\n" + cs.Bold("Available commands:"))
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			fmt.Printf("  %s\n", subCmd.Name())
		}
	}

	return nil
}
