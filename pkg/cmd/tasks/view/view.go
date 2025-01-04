package view

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/utils"
)

func NewCmdView(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View details of a specific task",
		Long: heredoc.Doc(`
				Display detailed information about a specific task,
				allowing you to analyze and manage it effectively.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return viewRun(f)
		},
	}

	return cmd
}

func viewRun(f factory.Factory) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	client, err := f.NewAsanaClient()
	if err != nil {
		return err
	}

	allTasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      cfg.Workspace.ID,
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
	taskNames := format.Tasks(allTasks)

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
		format.Date(task.DueOn), format.Projects(task.Projects))
	fmt.Println(format.Tags(task.Tags))
	fmt.Print(format.Notes(task.Notes))

	return nil
}
