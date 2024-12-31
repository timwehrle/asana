package list

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/api"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/utils"
)

func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List your tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun()
		},
	}

	return cmd
}

func listRun() error {
	token, err := auth.Get()
	if err != nil {
		return err
	}

	client := api.New(token)

	tasks, err := client.GetTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Println(utils.BoldUnderline().Sprintf("Your Tasks:"))
	for i, task := range tasks {
		fmt.Printf("%d. [%s] %s\n", i+1, utils.FormatDate(task.DueOn), task.Name)
	}

	return nil
}
