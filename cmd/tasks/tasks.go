package tasks

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/act/api"
	"github.com/timwehrle/act/internal/auth"
	"github.com/timwehrle/act/internal/prompter"
	"github.com/timwehrle/act/utils"
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

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return
		}

		taskNames := make([]string, len(tasks))
		for i, task := range tasks {
			taskNames[i] = fmt.Sprintf("[%s] %s", utils.FormatDate(task.DueOn), task.Name)
		}

		today := time.Now()

		selectMessage := fmt.Sprintf("Your Tasks (%s):", today.Format("Jan 02, 2006"))

		index, err := prompter.Select(selectMessage, taskNames)
		if err != nil {
			fmt.Println(err)
			return
		}

		selectedTask := tasks[index]

		fmt.Println("Selected Task:", selectedTask.Name)
	},
}
