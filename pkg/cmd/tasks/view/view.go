package view

import (
	"fmt"
	"time"

	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"

	"github.com/MakeNowJust/heredoc"
	taskshared "github.com/timwehrle/asana/pkg/cmd/tasks/shared"
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
	TaskID string
	Output string
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
				$ asana ts view
				$ asana tasks view --task 12001234 --output json`),
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
	cmd.Flags().StringVar(&opts.TaskID, "task", "", "Task GID to view directly without prompting")
	cmd.Flags().StringVar(&opts.TaskID, "task-id", "", "Task GID to view directly without prompting")
	cmd.Flags().StringVar(&opts.Output, "output", taskshared.OutputText, "Output format: text or json")

	return cmd
}

func viewRun(opts *ViewOptions) error {
	if err := taskshared.ValidateOutputMode("output", opts.Output); err != nil {
		return err
	}

	if opts.TaskID == "" {
		if err := taskshared.EnsureInteractiveAllowed(opts.IO, "--task <gid>"); err != nil {
			return err
		}
	}

	client, err := opts.Client()
	if err != nil {
		return err
	}

	if opts.TaskID != "" {
		task, err := taskshared.FetchTaskByID(client, opts.TaskID, []string{
			"name",
			"notes",
			"completed",
			"due_on",
			"permalink_url",
			"dependencies.name",
			"dependencies.completed",
			"custom_fields.name",
			"custom_fields.resource_subtype",
			"custom_fields.display_value",
			"custom_fields.text_value",
			"custom_fields.number_value",
			"custom_fields.boolean_value",
			"custom_fields.enum_value.name",
			"custom_fields.enum_value.color",
			"custom_fields.multi_enum_values.name",
			"custom_fields.multi_enum_values.color",
			"memberships.project.name",
			"memberships.section.name",
			"projects.name",
			"tags.name",
		})
		if err != nil {
			return err
		}
		return printTask(opts.IO, task, opts.Output)
	}

	cfg, err := opts.Config()
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

	return printTask(opts.IO, selectedTask, opts.Output)
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

func printTask(io *iostreams.IOStreams, task *asana.Task, output string) error {
	if taskshared.NormalizeOutputMode(output) == taskshared.OutputJSON {
		return taskshared.WriteJSON(io.Out, map[string]taskshared.TaskOutput{
			"task": taskshared.ToTaskOutput(task),
		})
	}

	cs := io.ColorScheme()

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
