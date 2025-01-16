package login

import (
	"fmt"

	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
)

type LoginOptions struct {
	factory.Factory
	IO     *iostreams.IOStreams
	Config struct {
		Workspace string
	}
}

func NewCmdLogin(f factory.Factory) *cobra.Command {
	opts := &LoginOptions{
		Factory: f,
		IO:      f.IOStreams(),
	}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to your Asana account",
		Long: heredoc.Docf(`
				Authenticate with Asana using a Personal Access Token (PAT).
				
				To get started:
				1. Visit https://app.asana.com/0/my-apps
				2. Click "Create new token"
				3. Give your token a description (e.g., "CLI Access")
				4. Copy the generated token`),
		Example: heredoc.Doc(`
					$ asana auth login
					$ asana auth login --workspace "Test Workspace"`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Config.Workspace, "workspace", "w", "", "Preselect the default workspace")

	return cmd
}

func runLogin(opts *LoginOptions) error {
	cs := opts.IO.ColorScheme()
	var token string

	token, err := auth.Get()
	if err == nil && token != "" {
		fmt.Fprintln(opts.IO.Out, "You are already logged in")
		return nil
	}

	fmt.Fprint(opts.IO.Out, heredoc.Doc(`
		Tip: You can generate a Personal Access Token here: https://app.asana.com/0/my-apps
	`))
	token, err = opts.Prompter().Token()
	if err != nil {
		return err
	}

	err = auth.ValidateToken(token)
	if err != nil {
		return err
	}

	client := asana.NewClientWithAccessToken(token)

	workspaces, err := client.AllWorkspaces()
	if err != nil {
		return err
	}

	err = auth.Set(token)
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, cs.SuccessIcon, "Logged in")

	if len(workspaces) == 0 {
		fmt.Fprintln(opts.IO.Out, "No workspaces found")
		return nil
	}

	var selectedWorkspace *asana.Workspace
	if opts.Config.Workspace != "" {
		for _, ws := range workspaces {
			if ws.ID == opts.Config.Workspace || ws.Name == opts.Config.Workspace {
				selectedWorkspace = ws
				break
			}
		}

		if selectedWorkspace == nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Workspace '%s' not found. Please select one from the list.\n", cs.ErrorIcon, opts.Config.Workspace)
		}
	}

	if selectedWorkspace == nil {
		names := make([]string, len(workspaces))
		for i, ws := range workspaces {
			names[i] = ws.Name
		}

		index, err := opts.Prompter().Select("Select a default workspace:", names)
		if err != nil {
			return err
		}

		selectedWorkspace = workspaces[index]
	}

	user, err := client.CurrentUser()
	if err != nil {
		return err
	}

	cfg := &config.Config{
		Username:  user.Name,
		Workspace: selectedWorkspace,
	}

	err = cfg.Save()
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.Out, "%s Default workspace set to %s\n", cs.SuccessIcon, cs.Bold(selectedWorkspace.Name))

	return nil
}
