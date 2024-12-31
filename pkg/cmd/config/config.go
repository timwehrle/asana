package config

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/config/get"
	"github.com/timwehrle/asana/pkg/cmd/config/set"
)

func NewCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure the Asana CLI config",
	}

	cmd.AddCommand(set.NewCmdSet())
	cmd.AddCommand(get.NewCmdGet())

	return cmd
}
