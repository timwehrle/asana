package update

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/cmd/tasks/shared"
	"github.com/timwehrle/asana/utils"
)

func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update details of a specific task",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateRun()
		},
	}

	return cmd
}

func updateRun() error {
	token, err := auth.Get()
	if err != nil {
		return err
	}

	defaultWorkspace, err := config.GetDefaultWorkspace()
	if err != nil {
		return err
	}

	client := asana.NewClientWithAccessToken(token)

	tasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      defaultWorkspace.ID,
		CompletedSince: "now",
	}, &asana.Options{
		Fields: []string{"due_on", "name"},
	})
	if err != nil {
		return err
	}

	taskNames := shared.FormatTasks(tasks)

	index, err := prompter.Select("Select the task to update", taskNames)

	selectedTask := tasks[index]

	err = selectedTask.Fetch(client)
	if err != nil {
		return err
	}

	actions := []string{"Mark as Completed", "Edit Task Name", "Cancel"}
	selectedAction, err := prompter.Select("What do you want to do with this task:", actions)

	switch selectedAction {
	case 0:
		return completeTask(client, selectedTask)
	case 1:
		return editTask(client, selectedTask)
	case 2:
		fmt.Println(utils.Success(), "Cancelled")
		return nil
	}

	return nil
}

func completeTask(client *asana.Client, task *asana.Task) error {
	completed := true
	taskBase := asana.TaskBase{
		Completed: &completed,
	}

	updateRequest := &asana.UpdateTaskRequest{
		TaskBase: taskBase,
	}

	err := task.Update(client, updateRequest)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	fmt.Println(utils.Success(), "Task completed")

	return nil
}

func editTask(client *asana.Client, task *asana.Task) error {
	input, err := prompter.Input("Enter the new task name:", "")
	if err != nil {
		return err
	}

	taskBase := asana.TaskBase{
		Name: input,
	}

	updateRequest := &asana.UpdateTaskRequest{
		TaskBase: taskBase,
	}

	err = task.Update(client, updateRequest)
	if err != nil {
		return err
	}

	fmt.Println(utils.Success(), "Task name updated")

	return nil
}
