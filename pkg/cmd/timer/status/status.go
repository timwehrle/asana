package status

import (
	"fmt"
	"sort"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type StatusOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
}

func NewCmdStatus(f factory.Factory, runF func(*StatusOptions) error) *cobra.Command {
	opts := &StatusOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the current tracked time of the selected task",
		Long: heredoc.Doc(`
				Show the total time tracked for a selected task in your Asana workspace.
			`),
		Example: heredoc.Doc(`
				# Show the tracked time of a selected task
				$ asana timer status
			`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF == nil {
				return runStatus(opts)
			}
			return runF(opts)
		},
	}

	return cmd
}

func runStatus(opts *StatusOptions) error {
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

	entries, _, err := task.GetTimeTrackingEntries(client, &asana.Options{
		Fields: []string{"created_by.name", "created_by.gid", "duration_minutes", "entered_on"},
	})
	if err != nil {
		return fmt.Errorf("failed to get time tracking entries: %w", err)
	}

	if len(entries) == 0 {
		io.Println("No time tracked yet for this task.")
		return nil
	}

	grouped := make(map[*asana.Date][]*asana.TimeTrackingEntry)
	totalMinutes := 0
	for _, entry := range entries {
		grouped[entry.EnteredOn] = append(grouped[entry.EnteredOn], entry)
		totalMinutes += entry.DurationMinutes
	}

	dates := make([]*asana.Date, 0, len(grouped))
	for d := range grouped {
		dates = append(dates, d)
	}

	sort.Slice(dates, func(i, j int) bool {
		return time.Time(*dates[i]).After(time.Time(*dates[j]))
	})

	io.Printf("\nTracked time entries on task %s:\n", cs.Bold(task.Name))
	for _, d := range dates {
		io.Printf("\n[%s]\n", format.Date(d))
		for _, entry := range grouped[d] {
			io.Printf("- %s tracked %s\n",
				entry.CreatedBy.Name,
				cs.Bold(format.Duration(entry.DurationMinutes)),
			)
		}
	}

	io.Printf("\nTotal tracked time: %s\n", cs.Bold(format.Duration(totalMinutes)))
	return nil
}

func selectTask(opts *StatusOptions, c *asana.Client) (*asana.Task, error) {
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
