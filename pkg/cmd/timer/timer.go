package timer

import (
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/cmd/timer/status"
	"github.com/timwehrle/asana/pkg/factory"
)

func NewCmdTimer(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timer",
		Short: "Manage the time tracking of your Asana tasks",
	}

	cmd.AddCommand(status.NewCmdStatus(f, nil))

	return cmd
}
