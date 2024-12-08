package status

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/timwehrle/act/api"
	"github.com/timwehrle/act/internal/auth"
	"github.com/timwehrle/act/internal/workspace"
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
				errCh <- fmt.Errorf("failed to fetch user information: %w", err)
			}
		}()

		wg.Wait()
		close(errCh)

		for err := range errCh {
			fmt.Println("Error:", err)
		}

		if tokenErr != nil {
			return
		}

		fmt.Println("API is operational.")
		fmt.Printf("Logged in as: %s (%s)\n", me.Name, me.GID)
		fmt.Printf("Default workspace: %s (%s)\n", name, gid)
	},
}
