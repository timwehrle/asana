package view

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/cmd/tasks/shared"
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

	defaultWorkspace, err := config.GetDefaultWorkspace()
	if err != nil {
		return err
	}

	client := asana.NewClientWithAccessToken(token)

	allTasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      defaultWorkspace.ID,
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
	taskNames := shared.FormatTasks(allTasks)

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
		utils.FormatDate(task.DueOn), shared.FormatProjects(task.Projects))
	fmt.Println(shared.FormatTags(task.Tags))
	fmt.Print(shared.FormatNotes(task.Notes))

	return nil
}
