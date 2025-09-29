package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type DeleteOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
}

func NewCmdDelete(f factory.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a time entry from a task",
		Long:  "Delete and remove a time entry from a selected Asana task.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF == nil {
				return runDelete(opts)
			}

			return runF(opts)
		},
	}

	return cmd
}

func runDelete(opts *DeleteOptions) error {
	io := opts.IO

	client, err := opts.Client()
	if err != nil {
		return err
	}

	task, err := selectTask(opts, client)
	if err != nil {
		return err
	}

	entries, _, err := task.GetTimeTrackingEntries(client, &asana.Options{
		Fields: []string{"created_by.name", "created_by.gid", "duration_minutes", "entered_on"},
	})
	if err != nil {
		return fmt.Errorf("failed to get time tracking entries: %w", err)
	}

	if len(entries) == 0 {
		io.Println("No time entries found for this task.")
		return nil
	}

	entryLabels := format.MapToStrings(entries, func(e *asana.TimeTrackingEntry) string {
		return fmt.Sprintf("%s — %s", e.CreatedBy.Name, format.Duration(e.DurationMinutes))
	})
	index, err := opts.Prompter.Select("Select a time entry to delete:", entryLabels)
	if err != nil {
		return fmt.Errorf("failed to select time entry: %w", err)
	}

	selectedEntry := entries[index]
	if err := selectedEntry.Delete(client); err != nil {
		return fmt.Errorf("failed to delete time tracking entry: %w", err)
	}

	cs := io.ColorScheme()
	io.Printf("%s Deleted time entry of %s created by %s\n", cs.SuccessIcon, format.Duration(selectedEntry.DurationMinutes), selectedEntry.CreatedBy.Name)

	return nil
}

func selectTask(opts *DeleteOptions, c *asana.Client) (*asana.Task, error) {
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
