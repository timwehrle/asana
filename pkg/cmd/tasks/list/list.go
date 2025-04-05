package list

import (
	"fmt"

	"github.com/timwehrle/asana/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
	"github.com/timwehrle/asana/pkg/sorting"
)

type SortOption string

const (
	SortAsc       SortOption = "asc"
	SortDesc      SortOption = "desc"
	SortDue       SortOption = "due"
	SortDueDesc   SortOption = "due-desc"
	SortCreatedAt SortOption = "created-at"
)

var validSortOptions = map[SortOption]struct{}{
	SortAsc:       {},
	SortDesc:      {},
	SortDue:       {},
	SortDueDesc:   {},
	SortCreatedAt: {},
}

type ListOptions struct {
	IO *iostreams.IOStreams

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)

	Sort  SortOption
	Limit int
}

func NewCmdList(f factory.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		Client: f.Client,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all tasks",
		Long: heredoc.Doc(`
				Retrieve and display a list of all tasks assigned to your Asana account.
				Tasks can be sorted by name, due date, or creation date.
			`),
		Example: heredoc.Doc(`
				# List all tasks
				$ asana tasks list

				# List tasks sorted by due date (descending)
				$ asana task list --sort due-desc
			`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateSortOption(opts.Sort); err != nil {
				return err
			}

			if runF != nil {
				return runF(opts)
			}

			return listRun(opts)
		},
	}

	cmd.Flags().
		StringVarP((*string)(&opts.Sort), "sort", "s", "", "Sort tasks by name, due date, creation date (options: asc, desc, due, due-desc, created-at)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "l", 0, "Limit the tasks to display")

	return cmd
}

func validateSortOption(opt SortOption) error {
	if opt == "" {
		return nil
	}

	if _, ok := validSortOptions[opt]; !ok {
		return fmt.Errorf(
			"invalid sort option %q. Available options: asc, desc, due, due-desc, created-at",
			opt,
		)
	}
	return nil
}

func listRun(opts *ListOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	tasks, err := fetchTasks(opts, cfg.Workspace.ID, opts.Limit)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return printEmptyMessage(opts.IO)
	}

	sortTasks(tasks, opts.Sort)

	return printTasks(opts.IO, cfg.Username, tasks)
}

func fetchTasks(opts *ListOptions, workspaceID string, limit int) ([]*asana.Task, error) {
	initialCapacity := 100
	if limit > 0 {
		initialCapacity = limit
	}

	client, err := opts.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create Asana client: %w", err)
	}

	query := &asana.TaskQuery{
		Assignee:       "me",
		Workspace:      workspaceID,
		CompletedSince: "now",
	}

	tasks := make([]*asana.Task, 0, initialCapacity)
	options := &asana.Options{
		Fields: []string{"name", "due_on", "created_at"},
		Limit:  limit,
	}

	for {
		batch, nextPage, err := client.QueryTasks(query, options)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tasks: %w", err)
		}

		tasks = append(tasks, batch...)

		if limit > 0 && len(tasks) > limit {
			tasks = tasks[:limit]
			break
		}

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	return tasks, nil
}

func sortTasks(tasks []*asana.Task, sortOption SortOption) {
	switch sortOption {
	case SortAsc:
		sorting.TaskSort.ByName(tasks)
	case SortDesc:
		sorting.TaskSort.ByNameDesc(tasks)
	case SortDue:
		sorting.TaskSort.ByDueDate(tasks)
	case SortDueDesc:
		sorting.TaskSort.ByDueDateDesc(tasks)
	case SortCreatedAt:
		sorting.TaskSort.ByCreatedAt(tasks)
	case "":
		// No sorting requested
	}
}

func printEmptyMessage(io *iostreams.IOStreams) error {
	fmt.Fprintln(io.Out, "No tasks found.")
	return nil
}

func printTasks(io *iostreams.IOStreams, username string, tasks []*asana.Task) error {
	cs := io.ColorScheme()

	fmt.Fprintf(io.Out, "\nTasks for %s:\n\n", cs.Bold(username))

	for i, task := range tasks {
		fmt.Fprintf(io.Out, "%d. [%s] %s\n",
			i+1,
			format.Date(task.DueOn),
			cs.Bold(task.Name),
		)
	}

	return nil
}
