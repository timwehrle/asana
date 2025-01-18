package config

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/config/get"
	"github.com/timwehrle/asana/pkg/cmd/config/set"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdConfig(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <subcommand>",
		Short: "Manage Asana CLI configuration",
		Long: heredoc.Doc(`
				Set and retrieve configuration settings for the Asana CLI tool.
		`),
	}

	cmd.AddCommand(set.NewCmdConfigSet(f, nil))
	cmd.AddCommand(get.NewCmdGet(f, nil))

	return cmd
}
