package view

import (
	"fmt"
	"time"

	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"

	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
)

type ViewOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
}

func NewCmdView(f factory.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "view",
		Short: "View details of a specific task",
		Example: heredoc.Doc(`
				$ asana tasks view
				$ asana ts view`),
		Long: heredoc.Doc(`
				Display detailed information about a specific task, allowing you to
				analyze and manage it effectively.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return viewRun(opts)
		},
	}

	return cmd
}

func viewRun(opts *ViewOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	client, err := opts.Client()
	if err != nil {
		return err
	}

	allTasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      cfg.Workspace.ID,
		CompletedSince: "now",
	}, &asana.Options{
		Fields: []string{"due_on", "name"},
	})
	if err != nil {
		return err
	}

	selectedTask, err := prompt(allTasks, opts.Prompter)
	if err != nil {
		return err
	}

	err = displayDetails(client, selectedTask, opts.IO)
	if err != nil {
		return err
	}

	return nil
}

func prompt(allTasks []*asana.Task, prompter prompter.Prompter) (*asana.Task, error) {
	taskNames := format.Tasks(allTasks)

	today := time.Now()
	selectMessage := fmt.Sprintf(
		"Your Tasks on %s (Select one for more details):",
		today.Format("Jan 02, 2006"),
	)

	index, err := prompter.Select(selectMessage, taskNames)
	if err != nil {
		return nil, err
	}

	return allTasks[index], nil
}

func displayDetails(client *asana.Client, task *asana.Task, io *iostreams.IOStreams) error {
	cs := io.ColorScheme()

	err := task.Fetch(client)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		io.Out,
		"%s | Due: %s | %s\n",
		cs.Bold(task.Name),
		format.Date(task.DueOn),
		format.Projects(task.Projects),
	)
	fmt.Fprintf(io.Out, "%s\n", format.Tags(task.Tags))
	fmt.Fprintln(io.Out, format.Notes(task.Notes))

	return nil
}
