package create

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/convert"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
)

type CreateOptions struct {
	cmdutils.BaseOptions
}

func NewCmdCreate(f factory.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		BaseOptions: cmdutils.BaseOptions{
			IO:       f.IOStreams,
			Prompter: f.Prompter,
			Config:   f.Config,
			Client:   f.Client,
		},
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

	task, err := cmdutils.SelectTask(&opts.BaseOptions, client)
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
