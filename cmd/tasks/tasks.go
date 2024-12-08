package tasks

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/act/api"
	"github.com/timwehrle/act/internal/auth"
)

var TasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage tasks",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := auth.Get()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := api.New(token)

		tasks, err := client.GetTasks()
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Name)
		}
	},
}
