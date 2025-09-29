package cmdutils

import (
	"fmt"

	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type BaseOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter
	Config   func() (*config.Config, error)
	Client   func() (*asana.Client, error)
}

func SelectTask(opts *BaseOptions, c *asana.Client) (*asana.Task, error) {
	cfg, err := opts.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	tasks, _, err := c.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      cfg.Workspace.ID,
		CompletedSince: "now",
	}, &asana.Options{
		Fields: []string{"name", "due_on"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	if len(tasks) == 0 {
		opts.IO.Println("No tasks found.")
		return nil, nil
	}

	taskNames := format.Tasks(tasks)
	index, err := opts.Prompter.Select("Select a task to view tracked time:", taskNames)
	if err != nil {
		return nil, fmt.Errorf("failed to select task: %w", err)
	}

	selectedTask := tasks[index]
	if err := selectedTask.Fetch(c); err != nil {
		return nil, fmt.Errorf("failed to fetch task details: %w", err)
	}

	return selectedTask, nil
}
