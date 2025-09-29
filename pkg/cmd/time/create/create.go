package create

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/convert"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type CreateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
}

func NewCmdCreate(f factory.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new time entry for a task",
		Long:  "Create and log a new time entry on a selected Asana task.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF == nil {
				return runCreate(opts)
			}

			return runF(opts)
		},
	}

	return cmd
}

func runCreate(opts *CreateOptions) error {
	io := opts.IO
	cs := io.ColorScheme()

	client, err := opts.Client()
	if err != nil {
		return err
	}

	task, err := selectTask(opts, client)
	if err != nil {
		return err
	}

	minutes, err := promptDuration(opts)
	if err != nil {
		return err
	}

	date, err := promptDate(opts)
	if err != nil {
		return err
	}

	result, err := task.CreateTimeTrackingEntry(client, &asana.CreateTimeTrackingEntryRequest{
		DurationMinutes: minutes,
		EnteredOn:       date,
	})
	if err != nil {
		return fmt.Errorf("failed to create time tracking entry: %w", err)
	}

	io.Printf("%s Logged %s to %q (%s)\n", cs.SuccessIcon, format.Duration(result.DurationMinutes), task.Name, format.HumanDate(*result.CreatedAt))
	return nil
}

func promptDuration(opts *CreateOptions) (int, error) {
	minutesStr, err := opts.Prompter.Input("How many minutes do you want to log? (e.g., 30)", "")
	if err != nil {
		return 0, fmt.Errorf("failed to read duration: %w", err)
	}

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil || minutes <= 0 {
		return 0, fmt.Errorf("invalid duration: must be a positive number")
	}
	return minutes, nil
}

func promptDate(opts *CreateOptions) (*asana.Date, error) {
	input, err := opts.Prompter.Input("Enter date [YYYY-MM-DD] or leave empty for today:", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read date: %w", err)
	}

	if input == "" {
		today := asana.Date(time.Now())
		return &today, nil
	}

	return convert.ToDate(input, time.DateOnly)
}

func selectTask(opts *CreateOptions, c *asana.Client) (*asana.Task, error) {
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
	index, err := opts.Prompter.Select("Select the task to see the time of:", taskNames)
	if err != nil {
		return nil, fmt.Errorf("failed to select task: %w", err)
	}

	selectedTask := tasks[index]
	if err := selectedTask.Fetch(c); err != nil {
		return nil, fmt.Errorf("failed to fetch task details: %w", err)
	}

	return selectedTask, nil
}
