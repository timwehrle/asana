package list

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/utils"
)

func NewCmdList(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all tasks",
		Long: heredoc.Doc(`
			Retrieve and display a list of all tasks assigned to your Asana account.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun(f)
		},
	}

	return cmd
}

func listRun(f factory.Factory) error {
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

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Println(utils.BoldUnderline().Sprintf("Your Tasks:"))
	for i, task := range tasks {
		fmt.Printf("%d. [%s] %s\n", i+1, format.Date(task.DueOn), task.Name)
	}

	return nil
}
