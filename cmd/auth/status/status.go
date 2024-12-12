package status

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"sync"

	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/workspace"
)

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of the user.",
	Long: heredoc.Doc(`
			Get the status of the current user and API.

			This command displays the API's operational status, 
			the logged-in user's username, and the default workspace.
	`),
	Example: heredoc.Doc(`
			# Display status 
			$ act auth status
	`),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			wg       sync.WaitGroup
			errCh    = make(chan error, 2)
			gid      string
			name     string
			me       api.User
			tokenErr error
			apiErr   error
		)

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			gid, name, err = workspace.LoadDefaultWorkspace()
			if err != nil {
				errCh <- fmt.Errorf("failed to load default workspace: %w", err)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := auth.Get()
			if err != nil {
				tokenErr = fmt.Errorf("failed to get token: %w", err)
				errCh <- tokenErr
				return
			}

			client := api.New(token)
			me, err = client.GetMe()
			if err != nil {
				apiErr = fmt.Errorf("failed to fetch user information: %w", err)
				errCh <- apiErr
				return
			}
		}()

		wg.Wait()
		close(errCh)

		if tokenErr != nil {
			fmt.Println("You are not logged in.")
			return
		}

		if apiErr != nil {
			fmt.Println("API is not operational.")
			return
		} else {
			fmt.Println("API is operational.")
		}

		fmt.Printf("Logged in as: %s (%s)\n", me.Name, me.GID)
		if gid == "" || name == "" {
			fmt.Println("No default workspace set.")
		} else {
			fmt.Printf("Default workspace: %s (%s)\n", name, gid)
		}
	},
}
