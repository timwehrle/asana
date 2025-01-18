package tasks

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/pkg/cmd/projects/shared"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type TasksOptions struct {
	factory.Factory
	IO *iostreams.IOStreams
}

func NewCmdTasks(f factory.Factory) *cobra.Command {
	opts := &TasksOptions{
		Factory: f,
		IO:      f.IOStreams(),
	}

	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "List tasks of a project",
		Long:  "Retrieve and display a list of all tasks under a project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTasks(opts)
		},
	}

	return cmd
}

func runTasks(opts *TasksOptions) error {
	cfg, err := opts.Factory.Config()
	if err != nil {
		return err
	}

	client, err := opts.Client()
	if err != nil {
		return err
	}

	var projects []*asana.Project
	workspace := &asana.Workspace{
		ID: cfg.Workspace.ID,
	}

	projects, err = shared.FetchAllProjects(client, workspace, 0)
	if err != nil {
		return err
	}

	projectNames := make([]string, 0, len(projects))
	for _, project := range projects {
		projectNames = append(projectNames, project.Name)
	}

	index, err := opts.Prompter().Select("Select a project:", projectNames)
	if err != nil {
		return fmt.Errorf("failed selecting a project: %w", err)
	}

	project := projects[index]

	tasks, _, err := project.Tasks(client)
	if err != nil {
		return fmt.Errorf("failed fetching tasks of project %s: %w", project.Name, err)
	}

	return printTasks(opts, project, tasks)
}

func printTasks(opts *TasksOptions, project *asana.Project, tasks []*asana.Task) error {
	cs := opts.IO.ColorScheme()

	fmt.Fprintf(opts.IO.Out, "\nTasks in %s:\n\n", cs.Bold(project.Name))

	if len(tasks) == 0 {
		fmt.Fprintf(opts.IO.Out, "No tasks found\n")
		return nil
	}

	for i, task := range tasks {
		fmt.Fprintf(opts.IO.Out, "%d. %s\n", i+1, cs.Bold(task.Name))
	}

	return nil
}
