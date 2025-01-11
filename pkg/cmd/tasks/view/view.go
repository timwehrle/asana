package view

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-go"
)

type ViewOptions struct {
	factory.Factory
	*iostreams.IOStreams
}

func NewCmdView(f factory.Factory) *cobra.Command {
	opts := &ViewOptions{
		Factory:   f,
		IOStreams: f.IOStreams(),
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
			return viewRun(opts)
		},
	}

	return cmd
}

func viewRun(opts *ViewOptions) error {
	cfg, err := opts.Factory.Config()
	if err != nil {
		return err
	}

	client, err := opts.NewAsanaClient()
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

	selectedTask, err := prompt(allTasks, opts.Factory)
	if err != nil {
		return err
	}

	err = displayDetails(client, selectedTask, opts)
	if err != nil {
		return err
	}

	return nil
}

func prompt(allTasks []*asana.Task, f factory.Factory) (*asana.Task, error) {
	taskNames := format.Tasks(allTasks)

	today := time.Now()
	selectMessage := fmt.Sprintf("Your Tasks on %s (Select one for more details):", today.Format("Jan 02, 2006"))

	index, err := f.Prompter().Select(selectMessage, taskNames)
	if err != nil {
		return nil, err
	}

	return allTasks[index], nil
}

func displayDetails(client *asana.Client, task *asana.Task, opts *ViewOptions) error {
	cs := opts.IOStreams.ColorScheme()

	err := task.Fetch(client)
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.Out, "%s | Due: %s | %s\n", cs.Bold(task.Name), format.Date(task.DueOn), format.Projects(task.Projects))
	fmt.Fprintf(opts.Out, "%s\n", format.Tags(task.Tags))
	fmt.Fprintln(opts.Out, format.Notes(task.Notes))

	return nil
}
