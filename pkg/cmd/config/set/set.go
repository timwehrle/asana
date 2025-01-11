package set

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type SetOptions struct {
	factory.Factory
	IO *iostreams.IOStreams
}

func NewCmdConfigSet(f factory.Factory) *cobra.Command {
	opts := &SetOptions{
		Factory: f,
		IO:      f.IOStreams(),
	}

	cmd := &cobra.Command{
		Use:   "set <key>",
		Short: "Update configuration with a value",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
				# Set a configuration value
				$ asana config set default-workspace
				$ asana config set dw
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSet(opts, args[0])
		},
	}

	return cmd
}

func runConfigSet(opts *SetOptions, key string) error {
	switch key {
	case "default-workspace", "dw":
		return setDefaultWorkspace(opts)
	default:
		return fmt.Errorf("unknown configuration key: %s. Available keys are: default-workspace (dw)", key)
	}
}

func setDefaultWorkspace(opts *SetOptions) error {
	cs := opts.IO.ColorScheme()

	client, err := opts.NewAsanaClient()
	if err != nil {
		return err
	}

	workspaces, err := client.AllWorkspaces()
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		fmt.Fprintln(opts.IO.Out, "No workspaces found")
		return nil
	}

	names := make([]string, len(workspaces))
	for i, ws := range workspaces {
		names[i] = ws.Name
	}

	index, err := opts.Prompter().Select("Select a new default workspace:", names)
	if err != nil {
		return fmt.Errorf("failed to select new workspace: %w", err)
	}

	selectedWorkspace := workspaces[index]

	cfg, err := opts.Factory.Config()
	if err != nil {
		return err
	}

	// Workspace must be uppercase here since the Set function works with the
	// interface names and workspace is uppercased.
	err = cfg.Set("Workspace", selectedWorkspace)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, "%s Default workspace set to %s\n", cs.SuccessIcon, cs.Bold(selectedWorkspace.Name))

	return nil
}
