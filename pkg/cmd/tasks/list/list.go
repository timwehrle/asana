package list

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/utils"
)

func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all tasks",
		Long: heredoc.Doc(`
			Retrieve and display a list of all tasks assigned to your Asana account.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun()
		},
	}

	return cmd
}

func listRun() error {
	token, err := auth.Get()
	if err != nil {
		return err
	}

	workspace, err := config.GetDefaultWorkspace()
	if err != nil {
		return err
	}

	client := asana.NewClientWithAccessToken(token)

	tasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      workspace.ID,
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
