package time

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/time/create"
	"github.com/timwehrle/asana/pkg/cmd/time/status"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdTimer(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "time",
		Short: "Manage time tracking for your Asana tasks",
		Long:  "Commands to track, delete and inspect time entries on your Asana tasks.",
	}

	cmd.AddCommand(status.NewCmdStatus(f, nil), create.NewCmdCreate(f, nil))

	return cmd
}
