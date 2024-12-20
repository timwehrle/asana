package tasks

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/prompter"
	"github.com/timwehrle/alfie/utils"
	"time"
)

func NewCmdTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tasksRun()
		},
	}

	return cmd
}

func tasksRun() error {
	token, err := auth.Get()
	if err != nil {
		return err
	}

	client := api.New(token)

	tasks, err := fetch(client)
	if err != nil {
		return err
	}

	selectedTask, err := prompt(tasks)
	if err != nil {
		return err
	}

	if err := displayDetails(client, selectedTask); err != nil {
		return err
	}

	if err := handleAction(client, selectedTask); err != nil {
		return err
	}

	return nil
}

func fetch(client *api.Client) ([]api.Task, error) {
	tasks, err := client.GetTasks()
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil, nil
	}

	return tasks, nil
}

func prompt(tasks []api.Task) (*api.Task, error) {
	taskNames := formatNames(tasks)

	today := time.Now()
	selectMessage := fmt.Sprintf("Your Tasks on %s (Select one for more details):", today.Format("Jan 02, 2006"))

	index, err := prompter.Select(selectMessage, taskNames)
	if err != nil {
		return nil, err
	}

	return &tasks[index], nil
}

func formatNames(tasks []api.Task) []string {
	taskNames := make([]string, len(tasks))
	for i, task := range tasks {
		taskNames[i] = fmt.Sprintf("[%s] %s", utils.FormatDate(task.DueOn), task.Name)
	}
	return taskNames
}

func displayDetails(client *api.Client, task *api.Task) error {
	detailedTask, err := client.GetTask(task.GID)
	if err != nil {
		return err
	}

	fmt.Printf("%s [%s], %s\n", utils.BoldUnderline.Sprint(detailedTask.Name),
		utils.FormatDate(detailedTask.DueOn), formatProjects(detailedTask.Projects))
	fmt.Println(formatTags(detailedTask.Tags))
	fmt.Print(formatNotes(detailedTask.Notes))

	return nil
}

func formatProjects(projects []api.Project) string {
	if len(projects) > 0 {
		projectNames := make([]string, len(projects))
		for i, project := range projects {
			projectNames[i] = project.Name
		}
		return "Projects: " + fmt.Sprintf("%s", projectNames)
	}
	return "Projects: None"
}

func formatTags(tags []api.Tag) string {
	if len(tags) > 0 {
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}
		return "Tags: " + fmt.Sprintf("%s", tagNames)
	}
	return "Tags: None"
}

func formatNotes(notes string) string {
	if notes != "" {
		return utils.BoldUnderline.Sprintf("Description:") + "\n" + notes + "\n"
	}
	return ""
}

func handleAction(client *api.Client, task *api.Task) error {
	actions := []string{"Mark as Done", "Delete Task", "Edit Task", "Cancel"}

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
