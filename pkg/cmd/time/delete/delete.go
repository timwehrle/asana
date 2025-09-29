package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
)

type DeleteOptions struct {
	cmdutils.BaseOptions
}

func NewCmdDelete(f factory.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		BaseOptions: cmdutils.BaseOptions{
			IO:       f.IOStreams,
			Prompter: f.Prompter,
			Config:   f.Config,
			Client:   f.Client,
		},
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

	task, err := cmdutils.SelectTask(&opts.BaseOptions, client)
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
