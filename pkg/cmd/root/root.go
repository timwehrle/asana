package root

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/cmd/tags"

	"github.com/spf13/cobra"
	service "github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/version"
	"github.com/timwehrle/asana/pkg/cmd/auth"
	"github.com/timwehrle/asana/pkg/cmd/config"
	"github.com/timwehrle/asana/pkg/cmd/projects"
	"github.com/timwehrle/asana/pkg/cmd/tasks"
	"github.com/timwehrle/asana/pkg/cmd/users"
	"github.com/timwehrle/asana/pkg/cmd/workspaces"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdRoot(f factory.Factory) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "asana <command> <subcommand> [flags]",
		Short:   "The Asana CLI tool",
		Version: version.Version,
		Long:    `Work with Asana from the command line.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip all checks for auth commands
			if isAuthCommand(os.Args) {
				return nil
			}

			// For non-auth commands, check authentication and load config
			err := service.Check()
			if err != nil {
				return err
			}

			return nil
		},
	}

	// Add auth command first
	cmd.AddCommand(auth.NewCmdAuth(f))

	// Only load config for non-auth commands
	if !isAuthCommand(os.Args) {
		cfg, err := f.Config()
		if err != nil {
			return nil, err
		}

		err = cfg.Load()
		if err != nil {
			return nil, err
		}

		err = cfg.Set("version", version.Version)
		if err != nil {
			return nil, err
		}
	}

	// Add other commands
	cmd.AddCommand(tasks.NewCmdTasks(f))
	cmd.AddCommand(projects.NewCmdProjects(f))
	cmd.AddCommand(workspaces.NewCmdWorkspace(f))
	cmd.AddCommand(users.NewCmdUsers(f))
	cmd.AddCommand(config.NewCmdConfig(f))
	cmd.AddCommand(tags.NewCmdTags(f))

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		showHelp(command, strings, os.Stdout)
	})

	cmd.SetUsageFunc(func(command *cobra.Command) error {
		return showRootUsage(command)
	})

	cmd.SetVersionTemplate(heredoc.Doc(`
	asana version {{ .Version }}
	https://github.com/timwehrle/asana/releases/tag/v{{ .Version }}
	`))

	return cmd, nil
}

// isAuthCommand checks if the command being run is an auth command
func isAuthCommand(args []string) bool {
	if len(args) < 2 {
		return false
	}
	return args[1] == "auth"
}
