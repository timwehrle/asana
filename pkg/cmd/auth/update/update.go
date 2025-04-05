package update

import (
	"fmt"

	"github.com/timwehrle/asana/internal/prompter"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type UpdateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter
}

func NewCmdUpdate(f factory.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
	}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update the Personal Access Token of your Asana account",
		Long:    "Update the current Personal Access Token of your Asana account.",
		Example: heredoc.Doc(`$ asana auth update`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runUpdate(opts)
		},
	}

	return cmd
}

func runUpdate(opts *UpdateOptions) error {
	cs := opts.IO.ColorScheme()

	newToken, err := opts.Prompter.Token()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	err = auth.ValidateToken(newToken)
	if err != nil {
		return err
	}

	err = auth.Set(newToken)
	if err != nil {
		return fmt.Errorf("failed to set new token: %w", err)
	}

	fmt.Fprintln(opts.IO.Out, cs.SuccessIcon, "Token updated")

	return nil
}
