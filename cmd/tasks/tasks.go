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

		task, err := client.GetTask(selectedTask.GID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s [%s], %s\n", utils.BoldUnderline.Sprint(task.Name), utils.FormatDate(task.DueOn), displayProjects(task.Projects))
		fmt.Println(displayTags(task.Tags))
		fmt.Print(displayNotes(task.Notes))
	},
}

func displayProjects(projects []api.Project) string {
	if len(projects) > 0 {
		var projectNames []string
		for _, project := range projects {
			projectNames = append(projectNames, project.Name)
		}
		return "Projects: " + fmt.Sprintf("%s", projectNames)
	}

	return "Projects: None"
}

func displayTags(tags []api.Tag) string {
	if len(tags) > 0 {
		var tagNames []string
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
		return "Tags: " + fmt.Sprintf("%s", tagNames)
	}
	return "Tags: None"
}

func displayNotes(notes string) string {
	if notes != "" {
		return fmt.Sprintf("\n%s\n", notes)
	}

	return ""
}
