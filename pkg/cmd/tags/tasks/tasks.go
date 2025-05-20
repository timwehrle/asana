package tasks

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type TasksOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter
	Config   func() (*config.Config, error)
	Client   func() (*asana.Client, error)

	ID string
}

func NewCmdTasks(f factory.Factory, runF func(*TasksOptions) error) *cobra.Command {
	opts := &TasksOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "List tasks with a given tag",
		Long: heredoc.Doc(`
				List all tasks associated with a selected tag in your Asana workspace.
			`),
		Example: heredoc.Doc(`
				# List all tasks with the selected tag
				$ asana tags tasks

				# List all tasks with a given tag
				$ asana tags tasks --id 1234
			`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(opts.ID) < 4 {
				return fmt.Errorf("ID should be longer than 4")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runTasks(opts)
		},
	}

	cmd.Flags().StringVar(&opts.ID, "id", "", "Specify a tag ID")

	return cmd
}

func runTasks(opts *TasksOptions) error {
	cs := opts.IO.ColorScheme()

	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := opts.Client()
	if err != nil {
		return fmt.Errorf("failed to initialize Asana client: %w", err)
	}

	var tag *asana.Tag
	if opts.ID == "" {
		tag, err = getTag(opts, cfg.Workspace.ID, client)
		if err != nil {
			return err
		}
	} else {
		tag = &asana.Tag{ID: opts.ID}
	}

	err = tag.Fetch(client)
	if err != nil {
		return fmt.Errorf("failed to fetch tag: %w", err)
	}

	tasks, _, err := tag.Tasks(client)
	if err != nil {
		return fmt.Errorf("failed to fetch tasks for tag %s: %w", tag.Name, err)
	}

	if len(tasks) == 0 {
		return fmt.Errorf("no tasks found for tag %s", tag.Name)
	}

	opts.IO.Printf("\nTasks with the tag %s:\n\n", cs.Bold(tag.Name))
	for i, task := range tasks {
		opts.IO.Printf("%2d. %s\n", i+1, cs.Bold(task.Name))
	}

	return nil
}

func getTag(opts *TasksOptions, workspaceID string, client *asana.Client) (*asana.Tag, error) {
	ws := &asana.Workspace{ID: workspaceID}

	tags, err := ws.AllTags(client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	names := format.MapToStrings(tags, func(tag *asana.Tag) string {
		return tag.Name
	})

	// TODO: Implement multiselection of tags since tasks can be assigned to more than one tag
	selected, err := opts.Prompter.Select("Select a tag: ", names)
	if err != nil {
		return nil, fmt.Errorf("failed to select a tag: %w", err)
	}
	tag := tags[selected]

	return tag, nil
}
