package cmd

import (
	"github.com/timwehrle/asana/pkg/cmd/tags"
	"os"

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

func NewCmdRoot() (*cobra.Command, error) {
	cmdFactory := factory.New()

	cmd := &cobra.Command{
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

	cmd.AddCommand(auth.NewCmdAuth(cmdFactory))
	cmd.AddCommand(tasks.NewCmdTasks(cmdFactory))
	cmd.AddCommand(projects.NewCmdProjects(cmdFactory))
	cmd.AddCommand(workspaces.NewCmdWorkspace(cmdFactory))
	cmd.AddCommand(users.NewCmdUsers(cmdFactory))
	cmd.AddCommand(config.NewCmdConfig(cmdFactory))
	cmd.AddCommand(tags.NewCmdTags(cmdFactory))

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		showHelp(command, strings, os.Stdout)
	})

	cmd.SetUsageFunc(func(command *cobra.Command) error {
		return showRootUsage(command)
	})

	return cmd, nil
}
