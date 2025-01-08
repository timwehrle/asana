package login

import (
	"fmt"

	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
)

func NewCmdLogin(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to your Asana account",
		Long: heredoc.Docf(`Authenticate with Asana using a Personal Access Token.
				Follow the steps in your Asana account to generate a token and use it
				with this command to enable CLI access.`),
		Example: heredoc.Doc(`
					$ asana auth login`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return loginRun(f)
		},
	}

	return cmd
}

func loginRun(f factory.Factory) error {
	var token string

	_, err := auth.Get()
	if err == nil {
		fmt.Println("You are already logged in.")
		return nil
	}

	fmt.Print(heredoc.Doc(`
		Tip: You can generate a Personal Access Token here: https://app.asana.com/0/my-apps
	`))
	token, err = f.Prompter().Token()
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

	fmt.Println(utils.Success(), "Logged in")

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
		return nil
	}

	names := make([]string, len(workspaces))
	for i, ws := range workspaces {
		names[i] = ws.Name
	}

	index, err := f.Prompter().Select("Select a default workspace:", names)
	if err != nil {
		return err
	}

	selectedWorkspace := workspaces[index]

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

	fmt.Printf("%s Default workspace set to %s\n", utils.Success(), utils.Bold().Sprint(selectedWorkspace.Name))

	return nil
}
