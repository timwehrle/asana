package status

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/timwehrle/alaric/api"
	"github.com/timwehrle/alaric/internal/auth"
	"github.com/timwehrle/alaric/internal/workspace"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current status of the user.",
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
			me, err = client.Me()
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
