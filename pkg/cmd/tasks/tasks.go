package tasks

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/prompter"
	"github.com/timwehrle/alfie/utils"
)

func NewCmdTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tasks",
		RunE: func(_ *cobra.Command, _ []string) error {
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

	err = handleAction(client, selectedTask)
	if err != nil {
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
	actions := []string{"Mark as Completed", "Edit Task Name", "Cancel"}

	selectedAction, err := prompter.Select("What would you like to do with this task?", actions)
	if err != nil {
		return err
	}

	switch selectedAction {
	case 0:
		return completeTask(client, task)
	case 1:
		return editTask(client, task)
	case 2:
		fmt.Println("Action cancelled.")
		return nil
	}

	return nil
}

func completeTask(client *api.Client, task *api.Task) error {
	if err := client.UpdateTask(task.GID, map[string]any{"completed": true}); err != nil {
		return err
	}
	fmt.Println("Task successfully marked as done.")

	return nil
}

func editTask(client *api.Client, task *api.Task) error {
	input, err := prompter.Input("What is the new name for the task?", "")
	if err != nil {
		return err
	}

	if err := client.UpdateTask(task.GID, map[string]any{"name": input}); err != nil {
		return err
	}
	fmt.Println("Task successfully edited.")
	return nil
}
