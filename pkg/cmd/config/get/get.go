package get

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type GetOptions struct {
	factory.Factory
	IO *iostreams.IOStreams
}

func NewCmdGet(f factory.Factory) *cobra.Command {
	opts := &GetOptions{
		Factory: f,
		IO:      f.IOStreams(),
	}

	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print the value of a given configuration key",
		Example: heredoc.Doc(`
				$ asana config get default-workspace
				$ asana config get dw`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigGet(opts, args[0])
		},
	}

	return cmd
}

func runConfigGet(opts *GetOptions, key string) error {
	cs := opts.IO.ColorScheme()

	switch key {
	case "default-workspace", "dw":
		cfg, err := opts.Factory.Config()
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.Out, "Default workspace is %s (%s)\n", cs.Bold(cfg.Workspace.Name), cfg.Workspace.ID)
		return nil

	default:
		return fmt.Errorf("unknown configuration key: %s. Available keys are: default-workspace (dw)", key)
	}
}
