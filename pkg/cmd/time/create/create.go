package create

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/convert"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
)

type CreateOptions struct {
	cmdutils.BaseOptions

	Minutes int
	DateStr string
	Date    *asana.Date
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
		Short: "Log time to a task",
		Long:  "Record a new time entry on a selected Asana task.",
		Example: heredoc.Doc(`
			# Log time via flags
			asana time create --minutes 30 --date 2025-01-06

			# Log time interactively
			asana time create --date 2025-01-06
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF == nil {
				return runCreate(opts)
			}

			return runF(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Minutes, "minutes", "m", 0, "Minutes to log (prompted if not set)")
	cmd.Flags().StringVar(&opts.DateStr, "date", "", "Entry date (YYYY-MM-DD, defaults to today)")

	return cmd
}

func (o *CreateOptions) Validate() error {
	if o.Minutes < 0 {
		return fmt.Errorf("minutes must be zero or a positive integer")
	}

	if o.DateStr != "" {
		date, err := convert.ToDate(o.DateStr, time.DateOnly)
		if err != nil {
			return fmt.Errorf("invalid date: %w", err)
		}
		o.Date = date
	} else {
		today := asana.Date(time.Now())
		o.Date = &today
	}
	return nil
}

func runCreate(opts *CreateOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	client, err := opts.Client()
	if err != nil {
		return err
	}

	task, err := cmdutils.SelectTask(&opts.BaseOptions, client)
	if err != nil {
		return err
	}

	var minutes int
	if opts.Minutes > 0 {
		minutes = opts.Minutes
	} else {
		minutes, err = promptDuration(opts)
		if err != nil {
			return err
		}
	}

	result, err := task.CreateTimeTrackingEntry(client, &asana.CreateTimeTrackingEntryRequest{
		DurationMinutes: minutes,
		EnteredOn:       opts.Date,
	})
	if err != nil {
		return fmt.Errorf("failed to create time tracking entry: %w", err)
	}

	opts.IO.Printf("%s Logged %s to %q on %s\n",
		opts.IO.ColorScheme().SuccessIcon,
		format.Duration(result.DurationMinutes),
		task.Name,
		format.Date(result.EnteredOn),
	)
	return nil
}

func promptDuration(opts *CreateOptions) (int, error) {
	minutesStr, err := opts.Prompter.Input("Enter minutes to log (e.g., 30):", "")
	if err != nil {
		return 0, fmt.Errorf("failed to read duration: %w", err)
	}

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil || minutes <= 0 {
		return 0, fmt.Errorf("invalid input: please enter a positive number")
	}
	return minutes, nil
}
