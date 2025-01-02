package view

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/internal/config"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/utils"
)

func NewCmdView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View details of a specific task",
		Long: heredoc.Doc(`
				Display detailed information about a specific task,
				allowing you to analyze and manage it effectively.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewRun()
		},
	}

	return cmd
}

func viewRun() error {
	token, err := auth.Get()
	if err != nil {
		return err
	}

	workspace, err := config.GetDefaultWorkspace()
	if err != nil {
		return err
	}

	client := asana.NewClientWithAccessToken(token)

	allTasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      workspace.ID,
		CompletedSince: "now",
	}, &asana.Options{
		Fields: []string{"due_on", "name"},
	})
	if err != nil {
		return err
	}

	selectedTask, err := prompt(allTasks)
	if err != nil {
		return err
	}

	err = displayDetails(client, selectedTask)
	if err != nil {
		return err
	}

	return nil
}

func prompt(allTasks []*asana.Task) (*asana.Task, error) {
	taskNames := FormatNames(allTasks)

	today := time.Now()
	selectMessage := fmt.Sprintf("Your Tasks on %s (Select one for more details):", today.Format("Jan 02, 2006"))

	index, err := prompter.Select(selectMessage, taskNames)
	if err != nil {
		return nil, err
	}

	return allTasks[index], nil
}

func displayDetails(client *asana.Client, task *asana.Task) error {
	err := task.Fetch(client)
	if err != nil {
		return err
	}

	fmt.Printf("%s | Due: %s | %s\n", utils.BoldUnderline().Sprint(task.Name),
		utils.FormatDate(task.DueOn), FormatProjects(task.Projects))
	fmt.Println(FormatTags(task.Tags))
	fmt.Print(FormatNotes(task.Notes))

	return nil
}

/*func handleAction(client *api.Client, task *api.Task) error {
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
*/
