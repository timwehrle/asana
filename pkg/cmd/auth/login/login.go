package login

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/timwehrle/asana/internal/prompter"

	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
)

type LoginOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)

	Workspace   string
	Token       string
	Interactive bool
}

func NewCmdLogin(f factory.Factory, runF func(*LoginOptions) error) *cobra.Command {
	opts := &LoginOptions{
		IO:       f.IOStreams,
		Config:   f.Config,
		Client:   f.Client,
		Prompter: f.Prompter,
	}

	var tokenStdin bool

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
					# Log in interactively and select a workspace
					$ asana auth login
					
					# Log in with a default workspace
					$ asana auth login --workspace "Test Workspace"
					
					# Log in with a token and set a default workspace
					$ asana auth login --workspace "Test Workspace" --with-token < mytoken.txt`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if tokenStdin {
				if opts.Workspace == "" {
					return fmt.Errorf(
						"workspace must be specified with --workspace when using --with-token",
					)
				}
				defer opts.IO.In.Close()
				token, err := io.ReadAll(opts.IO.In)
				if err != nil {
					return fmt.Errorf("failed to read token from standard input: %w", err)
				}
				opts.Token = strings.TrimSpace(string(token))
			}

			if opts.Token == "" {
				opts.Interactive = true
			}

			if runF != nil {
				return runF(opts)
			}

			return runLogin(opts)
		},
	}

	cmd.Flags().
		StringVarP(&opts.Workspace, "workspace", "w", "", "The default workspace to make calls to")
	cmd.Flags().BoolVar(&tokenStdin, "with-token", false, "Read token from standard input")

	return cmd
}

func runLogin(opts *LoginOptions) error {
	cs := opts.IO.ColorScheme()

	var token string
	token, err := auth.Get()
	if err == nil && token != "" {
		cfg := &config.Config{}
		if err := cfg.Load(); err != nil {
			if err := auth.Delete(); err != nil {
				return fmt.Errorf("failed to clear existing token: %w", err)
			}
		} else {
			fmt.Fprintln(opts.IO.Out, "You are already logged in")
			return nil
		}
	}

	if opts.Interactive {
		fmt.Fprint(opts.IO.Out, heredoc.Doc(`
		Tip: You can generate a Personal Access Token here: https://app.asana.com/0/my-apps
	`))
		token, err = opts.Prompter.Token()
		if err != nil {
			return err
		}
	} else {
		token = opts.Token
	}

	err = auth.ValidateToken(token)
	if err != nil {
		return err
	}

	client := asana.NewClientWithAccessToken(token)

	user, err := client.CurrentUser()
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

	var selectedWorkspace *asana.Workspace
	if opts.Workspace != "" {
		for _, ws := range workspaces {
			if ws.ID == opts.Workspace || strings.EqualFold(ws.Name, opts.Workspace) {
				selectedWorkspace = ws
				break
			}
		}

		if selectedWorkspace == nil {
			if !opts.Interactive {
				return fmt.Errorf(
					"%s Workspace '%s' not found. Please specify a valid workspace with --workspace",
					cs.ErrorIcon,
					opts.Workspace,
				)
			}

			fmt.Fprintf(
				opts.IO.ErrOut,
				"%s Workspace '%s' not found. Please select one from the list.\n",
				cs.ErrorIcon,
				opts.Workspace,
			)
		}
	}

	if selectedWorkspace == nil && opts.Interactive {
		names := make([]string, len(workspaces))
		for i, ws := range workspaces {
			names[i] = ws.Name
		}

		index, err := opts.Prompter.Select("Select a default workspace:", names)
		if err != nil {
			return err
		}

		selectedWorkspace = workspaces[index]
	}

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "asana-cli")
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	cfg := &config.Config{
		Username:  user.Name,
		Workspace: selectedWorkspace,
	}

	err = cfg.Save()
	if err != nil {
		return err
	}

	err = auth.Set(token)
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, cs.SuccessIcon, "Logged in")
	if selectedWorkspace != nil {
		fmt.Fprintf(
			opts.IO.Out,
			"%s Default workspace set to %s\n",
			cs.SuccessIcon,
			cs.Bold(selectedWorkspace.Name),
		)
	}

	return nil
}
