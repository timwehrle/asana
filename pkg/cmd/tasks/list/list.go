package list

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/sorting"
	"github.com/timwehrle/asana/utils"
)

const (
	sortAsc       = "asc"
	sortDesc      = "desc"
	sortDue       = "due"
	sortDueDesc   = "due-desc"
	sortCreatedAt = "created-at"
)

type options struct {
	Sort string
}

func NewCmdList(f factory.Factory) *cobra.Command {
	opts := &options{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all tasks",
		Long:    heredoc.Doc(`Retrieve and display a list of all tasks assigned to your Asana account.`),
		Example: heredoc.Doc(`
				$ asana tasks list
				$ asana tasks ls
				$ asana ts ls
			`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRun(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Sort, "sort", "s", "", "Sort tasks by name, due date, creation date (options: asc, desc, due, due-desc, created-at)")

	return cmd
}

func listRun(f factory.Factory, opts *options) error {
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
		Fields: []string{"due_on", "name", "created_at"},
	})
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	if err := applySorting(tasks, opts.Sort); err != nil {
		return err
	}

	fmt.Printf("\nTasks for %s:\n\n", utils.Bold().Sprint(cfg.Username))
	for i, task := range tasks {
		fmt.Printf("%d. [%s] %s\n", i+1, format.Date(task.DueOn), utils.Bold().Sprint(task.Name))
	}

	return nil
}

func applySorting(tasks []*asana.Task, sortOption string) error {
	switch sortOption {
	case sortAsc:
		sorting.TaskSort.ByName(tasks)
	case sortDesc:
		sorting.TaskSort.ByNameDesc(tasks)
	case sortDue:
		sorting.TaskSort.ByDueDate(tasks)
	case sortDueDesc:
		sorting.TaskSort.ByDueDateDesc(tasks)
	case sortCreatedAt:
		sorting.TaskSort.ByCreatedAt(tasks)
	case "":
		// No sorting
	default:
		return errors.New("invalid sort option. Available options: asc, desc, due, due-desc, created-at")
	}

	return nil
}
