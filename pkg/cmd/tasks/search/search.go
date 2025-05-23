package search

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
	"strings"
)

type SearchOptions struct {
	IO     *iostreams.IOStreams
	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)

	Query           string
	Type            string
	Assignee        []string
	ExcludeAssignee []string
	TagsAll         []string
	SortAscending   bool
	CreatorAny      []string
	ExcludeCreator  []string
	Blocked         bool
	SortBy          string
	DueOnBefore     string
	DueOnAfter      string
	DueOn           string
	DueAtBefore     string
	DueAtAfter      string
}

func (o *SearchOptions) join(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	return strings.Join(ss, ",")
}

var validSortBy = []string{
	"due_date",
	"created_at",
	"completed_at",
	"likes",
	"modified_at",
}

func NewCmdSearch(f factory.Factory, runF func(*SearchOptions) error) *cobra.Command {
	opts := &SearchOptions{
		IO:     f.IOStreams,
		Config: f.Config,
		Client: f.Client,
	}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for tasks in your workspace",
		Long: heredoc.Doc(`
					Search for tasks in your Asana workspace with various filters and sorting options.

					This command allows you to search for tasks by text, assignee, creator, tags and more.
					Results can be sorted according to your preference.
				`),
		Example: heredoc.Doc(`
					# Search for milestone tasks assigned to you
					$ asana tasks search --type milestone --assignee me --sort-asc
		
					# Search for tasks containing "UI refresh" not assigned to specific users
					$ asana tasks search --query "UI refresh" --exclude-assignee 1234,5678 --tags-all 1234,4567
				`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutils.ValidateStringEnum("sort-by", opts.SortBy, validSortBy); err != nil {
				return err
			}
			if err := cmdutils.ValidateDate("due-on", opts.DueOn); err != nil {
				return err
			}
			if err := cmdutils.ValidateDate("due-on-before", opts.DueOnBefore); err != nil {
				return err
			}
			if err := cmdutils.ValidateDate("due-on-after", opts.DueOnAfter); err != nil {
				return err
			}
			if err := cmdutils.ValidateDate("due-at-before", opts.DueAtBefore); err != nil {
				return err
			}
			if err := cmdutils.ValidateDate("due-at-after", opts.DueAtAfter); err != nil {
				return err
			}
			if opts.DueOn != "" && (opts.DueOnBefore != "" || opts.DueOnAfter != "") {
				return fmt.Errorf("--due-on cannot be used with --due-on-before or --due-on-after")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF == nil {
				return runSearch(opts)
			}
			return runF(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Query, "query", "q", "", "Perform full-text search on task names and descriptions")
	cmd.Flags().StringVar(&opts.Type, "type", "default_task", "Resource subtype to filter tasks (e.g., default_task, milestone)")
	cmd.Flags().StringSliceVarP(&opts.Assignee, "assignee", "a", []string{"me"}, "Comma-separated list of assignee user IDs (e.g., 1234,me)")
	cmd.Flags().StringSliceVar(&opts.ExcludeAssignee, "exclude-assignee", nil, "Comma separated list of user IDs to exclude from the search (e.g., 1234,5678)")
	cmd.Flags().StringSliceVar(&opts.TagsAll, "tags-all", nil, "Comma-separated list of tags to include in the search")
	cmd.Flags().BoolVar(&opts.SortAscending, "sort-asc", false, "Sort results in ascending order")
	cmd.Flags().StringVar(&opts.SortBy, "sort-by", "modified_at", "Sort results by one of: due_date, created_at, completed_at, likes or modified_at")
	cmd.Flags().StringSliceVar(&opts.CreatorAny, "creator-any", nil, "Comma-separated list of user IDs to include in the search")
	cmd.Flags().StringSliceVar(&opts.ExcludeCreator, "exclude-creator", nil, "Comma-separated list of user IDs to exclude from the search")
	cmd.Flags().BoolVar(&opts.Blocked, "is-blocked", false, "Filter to tasks with incomplete dependencies")
	cmd.Flags().StringVar(&opts.DueOnBefore, "due-on-before", "", "Filter to tasks due before a specific date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.DueOnAfter, "due-on-after", "", "Filter to tasks due after a specific date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.DueOn, "due-on", "", "Filter to tasks due on a specific date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.DueAtBefore, "due-at-before", "", "Filter to tasks due at or before a specific date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.DueAtAfter, "due-at-after", "", "Filter to tasks due at or after a specific date (YYYY-MM-DD)")

	return cmd
}

func runSearch(opts *SearchOptions) error {
	io := opts.IO
	cs := io.ColorScheme()

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	client, err := opts.Client()
	if err != nil {
		return err
	}

	workspace := cfg.Workspace

	query := &asana.SearchTasksQuery{
		Text:            opts.Query,
		SortBy:          opts.SortBy,
		ResourceSubtype: opts.Type,
		AssigneeAny:     opts.join(opts.Assignee),
		AssigneeNot:     opts.join(opts.ExcludeAssignee),
		TagsAll:         opts.join(opts.TagsAll),
		SortAscending:   opts.SortAscending,
		CreatedByAny:    opts.join(opts.CreatorAny),
		CreatedByNot:    opts.join(opts.ExcludeCreator),
		IsBlocked:       opts.Blocked,
		DueOnBefore:     opts.DueOnBefore,
		DueOnAfter:      opts.DueOnAfter,
		DueOn:           opts.DueOn,
		DueAtBefore:     opts.DueAtBefore,
		DueAtAfter:      opts.DueAtAfter,
	}

	options := &asana.Options{
		Fields: []string{"name", "due_on"},
	}

	tasks, err := workspace.SearchTasks(client, query, options)
	if err != nil {
		return fmt.Errorf("failed searching tasks: %w", err)
	}

	if len(tasks) == 0 {
		io.Println("No tasks found matching your criteria.")
		io.Println("- Try broadening your search by removing some filters")
		io.Println("- Check if the assignee or creator IDs are correct")
		io.Println("- If searching by text, try using fewer or different keywords")
		if opts.Type != "default_task" {
			io.Println(fmt.Sprintf("- Try changing the task type from '%s' to 'default_task'", opts.Type))
		}
		return nil
	}

	io.Printf("\nTasks assigned to %s:\n\n", cs.Bold(strings.Join(opts.Assignee, ", ")))

	for i, task := range tasks {
		io.Printf("%2d. [%s] %s\n", i+1, format.Date(task.DueOn), task.Name)
	}

	return nil
}
