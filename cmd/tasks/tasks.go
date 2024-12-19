package tasks

import (
	"fmt"
	"github.com/timwehrle/alfie/pkg"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/internal/auth"
	"github.com/timwehrle/alfie/internal/prompter"
	"github.com/timwehrle/alfie/utils"
)

var Cmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.Get()
		if err != nil {
			return err
		}

		client := api.New(token)

		tasks, err := fetchTasks(client)
		if err != nil {
			return err
		}

		selectedTask, err := promptForTaskSelection(tasks)
		if err != nil {
			return err
		}

		if err := displayTaskDetails(client, selectedTask); err != nil {
			return err
		}

		if err := pkg.HandleTaskAction(client, selectedTask); err != nil {
			return err
		}

		return nil
	},
}

func fetchTasks(client *api.Client) ([]api.Task, error) {
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

func promptForTaskSelection(tasks []api.Task) (*api.Task, error) {
	taskNames := formatTaskNames(tasks)

	today := time.Now()
	selectMessage := fmt.Sprintf("Your Tasks on %s (Select one for more details):", today.Format("Jan 02, 2006"))

	index, err := prompter.Select(selectMessage, taskNames)
	if err != nil {
		return nil, err
	}

	return &tasks[index], nil
}

func formatTaskNames(tasks []api.Task) []string {
	taskNames := make([]string, len(tasks))
	for i, task := range tasks {
		taskNames[i] = fmt.Sprintf("[%s] %s", utils.FormatDate(task.DueOn), task.Name)
	}
	return taskNames
}

func displayTaskDetails(client *api.Client, task *api.Task) error {
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
