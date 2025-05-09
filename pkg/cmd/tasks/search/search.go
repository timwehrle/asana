package search

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/factory"
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
		Example: heredoc.Doc(`
					$ asana tasks search --type milestone --assignee me --sort-asc

					$ asana tasks search --query "UI refresh" --exclude-assignee 1234,5678 --tags-all 1234,4567
				`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cmdutils.ValidateStringEnum("sort-by", opts.SortBy, validSortBy)
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
	}

	tasks, err := workspace.SearchTasks(client, query)
	if err != nil {
		return fmt.Errorf("failed searching tasks: %w", err)
	}

	if len(tasks) == 0 {
		io.Println("No tasks found")
		return nil
	}

	io.Printf("\nTasks assigned to %s:\n\n", cs.Bold(strings.Join(opts.Assignee, ", ")))

	for i, task := range tasks {
		io.Printf("%2d. %s\n", i+1, task.Name)
	}

	return nil
}
