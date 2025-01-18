package get

import (
	"fmt"
	"github.com/timwehrle/asana/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type GetOptions struct {
	Config func() (*config.Config, error)
	IO     *iostreams.IOStreams
}

func NewCmdGet(f factory.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		Config: f.Config,
		IO:     f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print the value of a given configuration key",
		Example: heredoc.Doc(`
				$ asana config get default-workspace
				$ asana config get dw`),
		ValidArgs: []string{"default-workspace", "dw"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runConfigGet(opts, args[0])
		},
	}

	return cmd
}

func runConfigGet(opts *GetOptions, key string) error {
	cs := opts.IO.ColorScheme()

	switch key {
	case "default-workspace", "dw":
		cfg, err := opts.Config()
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.Out, "Default workspace is %s (%s)\n", cs.Bold(cfg.Workspace.Name), cfg.Workspace.ID)
	}

	return nil
}
