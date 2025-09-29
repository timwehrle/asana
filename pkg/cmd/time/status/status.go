package status

import (
	"fmt"
	"sort"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
)

type StatusOptions struct {
	cmdutils.BaseOptions
}

func NewCmdStatus(f factory.Factory, runF func(*StatusOptions) error) *cobra.Command {
	opts := &StatusOptions{
		BaseOptions: cmdutils.BaseOptions{
			IO:       f.IOStreams,
			Prompter: f.Prompter,
			Config:   f.Config,
			Client:   f.Client,
		},
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show tracked time for a task",
		Long: heredoc.Doc(`
				Display all time entries logged on a selected Asana task, grouped by date,
				along with the total tracked time.
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

type GroupedEntries struct {
	Date    time.Time
	Label   string
	Entries []*asana.TimeTrackingEntry
}

func runStatus(opts *StatusOptions) error {
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

	groups, total, err := groupEntries(entries)
	if err != nil {
		return err
	}

	io.Printf("\nTime entries for task: %s\n", cs.Bold(task.Name))
	for _, g := range groups {
		io.Printf("\n[%s]\n", g.Label)
		for _, entry := range g.Entries {
			io.Printf(" • %s — %s\n",
				entry.CreatedBy.Name,
				cs.Bold(format.Duration(entry.DurationMinutes)),
			)
		}
	}

	io.Printf("\nTotal: %s\n", cs.Bold(format.Duration(total)))
	return nil
}

func groupEntries(entries []*asana.TimeTrackingEntry) ([]GroupedEntries, int, error) {
	m := map[string]*GroupedEntries{}
	total := 0

	for _, e := range entries {
		if e.EnteredOn == nil {
			continue
		}

		key := time.Time(*e.EnteredOn).Format(time.DateOnly)
		t, err := time.Parse(time.DateOnly, key)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid entered_on date: %w", err)
		}

		if _, ok := m[key]; !ok {
			m[key] = &GroupedEntries{
				Date:  t,
				Label: format.HumanDate(t),
			}
		}

		m[key].Entries = append(m[key].Entries, e)
		total += e.DurationMinutes
	}

	groups := make([]GroupedEntries, 0, len(m))
	for _, g := range m {
		groups = append(groups, *g)
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Date.After(groups[j].Date)
	})

	return groups, total, nil
}
