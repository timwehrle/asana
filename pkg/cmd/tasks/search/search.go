package search

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
	"strings"
)

type SearchOptions struct {
	IO     *iostreams.IOStreams
	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)

	Text            string
	ResourceSubtype string
	Assignee        []string
	AssigneeNot     []string
	TagsAll         []string
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF == nil {
				return runSearch(opts)
			}
			return runF(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Text, "text", "t", "", "Perform full-text search on task names and descriptions")
	cmd.Flags().StringVar(&opts.ResourceSubtype, "resource-subtype", "default_task", "Resource subtype to filter tasks (e.g., default_task, milestone)")
	cmd.Flags().StringSliceVarP(&opts.Assignee, "assignee", "a", []string{"me"}, "Comma-separated list of assignee user IDs (e.g., 1234,me)")
	cmd.Flags().StringSliceVar(&opts.AssigneeNot, "not-assignee", nil, "Comma separated list of user IDs to exclude from the search (e.g., 1234,5678)")
	cmd.Flags().StringSliceVar(&opts.TagsAll, "tags", nil, "Comma-separated list of tags to include in the search")

	return cmd
}

func runSearch(opts *SearchOptions) error {
	cs := opts.IO.ColorScheme()
	ioS := opts.IO

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
		Text:            opts.Text,
		SortBy:          "modified_at",
		ResourceSubtype: opts.ResourceSubtype,
		AssigneeAny:     strings.Join(opts.Assignee, ","),
		AssigneeNot:     strings.Join(opts.AssigneeNot, ","),
		TagsAll:         strings.Join(opts.TagsAll, ","),
		SortAscending:   false,
	}

	tasks, err := workspace.SearchTasks(client, query)
	if err != nil {
		return fmt.Errorf("failed searching tasks: %w", err)
	}

	if len(tasks) == 0 {
		ioS.Println("No tasks found")
		return nil
	}

	ioS.Printf("\nTasks assigned to %s:\n\n", cs.Bold(strings.Join(opts.Assignee, ", ")))

	for i, task := range tasks {
		ioS.Printf("%2d. %s\n", i+1, task.Name)
	}

	return nil
}
