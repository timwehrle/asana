package pkg

import (
	"fmt"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/prompter"
)

func HandleTaskAction(client *api.Client, task *api.Task, additionalActions ...string) error {
	actions := []string{"Mark as Done", "Delete Task", "Edit Task", "Cancel"}
	actions = append(actions, additionalActions...)

	selectedAction, err := prompter.Select("What would you like to do with this task?", actions)
	if err != nil {
		return err
	}

	switch selectedAction {
	case 0:
		return markTaskAsDone(client, task)
	case 3:
		fmt.Println("Action cancelled.")
		return nil
	}

	return nil
}

func markTaskAsDone(client *api.Client, task *api.Task) error {
	confirm, err := prompter.Confirm("Do you want to mark the task as done?", "No")
	if err != nil {
		return err
	}

	if confirm {
		if err := client.MarkTaskAsDone(task.GID); err != nil {
			return err
		}
		fmt.Println("Task successfully marked as done.")
	} else {
		fmt.Println("Task not marked as done.")
	}

	return nil
}
