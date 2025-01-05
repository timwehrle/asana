package update

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/convert"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/utils"
	"strings"
	"time"
)

func NewCmdUpdate(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update details of a specific task",
		Long:  "Retrieve task details and select one for updating it.",
		Example: heredoc.Doc(`
			$ asana tasks update
			$ asana ts update`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateRun(f)
		},
	}

	return cmd
}

func updateRun(f factory.Factory) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	client, err := f.NewAsanaClient()
	if err != nil {
		return err
	}

	tasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      cfg.Workspace.ID,
		CompletedSince: "now",
	}, &asana.Options{
		Fields: []string{"due_on", "name"},
	})
	if err != nil {
		return err
	}

	taskNames := format.Tasks(tasks)

	index, err := prompter.Select("Select the task to update", taskNames)
	if err != nil {
		return err
	}

	selectedTask := tasks[index]

	err = selectedTask.Fetch(client)
	if err != nil {
		return err
	}

	actions := []string{
		"Mark as Completed",
		"Edit Task Name",
		"Edit Description",
		"Set Due Date",
		"Cancel",
	}
	selectedAction, err := prompter.Select("What do you want to do with this task:", actions)
	if err != nil {
		return fmt.Errorf("failed selecting an action: %w", err)
	}

	switch selectedAction {
	case 0:
		return completeTask(client, selectedTask)
	case 1:
		return editTask(client, selectedTask)
	case 2:
		return editDescription(client, selectedTask)
	case 3:
		return setDueDate(client, selectedTask)
	case 4:
		fmt.Println(utils.Success(), "Operation canceled. You can rerun the command to try again.")
		return nil
	}

	return nil
}

func setDueDate(client *asana.Client, task *asana.Task) error {
	input, err := prompter.Input("Enter the new due date (YYYY-MM-DD):", "")
	if err != nil {
		return err
	}

	dueDate, err := convert.ToDate(input, time.DateOnly)
	if err != nil {
		return fmt.Errorf("failed parsing the date: %w", err)
	}

	updateRequest := &asana.UpdateTaskRequest{
		TaskBase: asana.TaskBase{
			DueOn: dueDate,
		},
	}

	err = task.Update(client, updateRequest)
	if err != nil {
		return fmt.Errorf("failed updating task due date: %w", err)
	}

	fmt.Println(utils.Success(), "Due date updated")

	return nil
}

func editDescription(client *asana.Client, task *asana.Task) error {
	existingDescription := strings.TrimSpace(task.Notes)

	newDescription, err := prompter.Editor("Edit the description:", existingDescription)
	if err != nil {
		return err
	}
	newDescription = strings.TrimSpace(newDescription)

	if newDescription != existingDescription {
		updateRequest := &asana.UpdateTaskRequest{
			TaskBase: asana.TaskBase{
				Notes: newDescription,
			},
		}

		err = task.Update(client, updateRequest)
		if err != nil {
			return fmt.Errorf("failed to update task notes: %w", err)
		}

		fmt.Println(utils.Success(), "Description updated")
	} else {
		fmt.Println("No changes made to description")
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
		return fmt.Errorf("failed to update task completion: %w", err)
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
