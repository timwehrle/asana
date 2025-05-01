package search

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type SearchOptions struct {
	IO     *iostreams.IOStreams
	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
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

	return cmd
}

func runSearch(opts *SearchOptions) error {
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
		SortBy:          "modified_at",
		ResourceSubtype: "default_task",
		AssigneeAny:     "me",
		SortAscending:   false,
	}

	tasks, err := workspace.SearchTasks(client, query)
	if err != nil {
		return fmt.Errorf("failed searching tasks: %w", err)
	}

	if len(tasks) == 0 {
		opts.IO.Println("No tasks found")
	}

	for i, task := range tasks {
		opts.IO.Printf("%2d. %s\n", i+1, task.Name)
	}

	return nil
}
