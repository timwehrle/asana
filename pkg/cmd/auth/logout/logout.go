package logout

import (
	"fmt"

	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
)

type LogoutOptions struct {
	IO *iostreams.IOStreams
}

func NewCmdLogout(f factory.Factory, runF func(options *LogoutOptions) error) *cobra.Command {
	opts := &LogoutOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of your Asana account",
		Long: heredoc.Doc(`
				Log out of your current Asana account by removing locally
				stored credentials.

				This action revokes CLI access to the Asana API.`),
		Example: heredoc.Doc(`$ asana auth logout`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runLogout(opts)
		},
	}

	return cmd
}

func runLogout(opts *LogoutOptions) error {
	cs := opts.IO.ColorScheme()

	err := auth.Check()
	if err != nil {
		return err
	}

	err = auth.Delete()
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, cs.SuccessIcon, "Logged out")

	return nil
}
