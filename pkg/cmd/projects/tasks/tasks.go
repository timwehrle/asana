package tasks

import (
	"fmt"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/pkg/cmd/projects/shared"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type TasksOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)

	WithSections bool
}

type sectionTasks struct {
	section *asana.Section
	tasks   []*asana.Task
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
		Short: "List tasks of a project",
		Long:  "Retrieve and display a list of all tasks under a project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runTasks(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.WithSections, "sections", "s", false, "Group tasks by sections")

	return cmd
}

func runTasks(opts *TasksOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	client, err := opts.Client()
	if err != nil {
		return err
	}

	var project *asana.Project
	project, err = selectProject(opts, client, cfg.Workspace.ID)

	if opts.WithSections {
		return listTasksWithSections(opts, client, project)
	}

	return listAllTasks(opts, client, project)
}

func selectProject(opts *TasksOptions, client *asana.Client, workspaceID string) (*asana.Project, error) {
	projects, err := shared.FetchAllProjects(client, &asana.Workspace{ID: workspaceID}, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Fprintln(opts.IO.Out, "No projects found")
		return nil, nil
	}

	projectNames := make([]string, len(projects))
	for i, project := range projects {
		projectNames[i] = project.Name
	}

	index, err := opts.Prompter.Select("Select a project:", projectNames)
	if err != nil {
		return nil, fmt.Errorf("failed to select a project: %w", err)
	}

	return projects[index], nil
}

func listAllTasks(opts *TasksOptions, client *asana.Client, project *asana.Project) error {
	tasks := make([]*asana.Task, 0, 100)
	options := &asana.Options{}

	for {
		batch, nextPage, err := project.Tasks(client, options)
		if err != nil {
			return fmt.Errorf("failed to fetch tasks for project %q: %w", project.Name, err)
		}

		tasks = append(tasks, batch...)

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	return displayTasks(opts, project, tasks)
}

func listTasksWithSections(opts *TasksOptions, client *asana.Client, project *asana.Project) error {
	sections := make([]*asana.Section, 0, 20)
	options := &asana.Options{}

	for {
		batch, nextPage, err := project.Sections(client, options)
		if err != nil {
			return err
		}

		sections = append(sections, batch...)

		if nextPage == nil || nextPage.Offset == "" {
			break
		}

		options.Offset = nextPage.Offset
	}

	sectionsWithTasks := make([]sectionTasks, 0, len(sections))

	for _, section := range sections {
		tasks := make([]*asana.Task, 0, 50)
		options := &asana.Options{}

		for {
			batch, nextPage, err := section.Tasks(client, options)
			if err != nil {
				return fmt.Errorf("failed to fetch tasks for section %q: %w", section.Name, err)
			}

			tasks = append(tasks, batch...)

			if nextPage == nil || nextPage.Offset == "" {
				break
			}

			options.Offset = nextPage.Offset
		}

		sectionsWithTasks = append(sectionsWithTasks, sectionTasks{
			section: section,
			tasks:   tasks,
		})
	}

	return displayTasksBySection(opts, project, sectionsWithTasks)
}

func displayTasks(opts *TasksOptions, project *asana.Project, tasks []*asana.Task) error {
	cs := opts.IO.ColorScheme()
	out := opts.IO.Out

	fmt.Fprintf(out, "\nTasks in %s:\n\n", cs.Bold(project.Name))

	if len(tasks) == 0 {
		fmt.Fprintf(opts.IO.Out, "No tasks found\n")
		return nil
	}

	for i, task := range tasks {
		fmt.Fprintf(out, "%d. %s\n", i+1, cs.Bold(task.Name))
	}

	return nil

}

func displayTasksBySection(opts *TasksOptions, project *asana.Project, sections []sectionTasks) error {
	cs := opts.IO.ColorScheme()
	out := opts.IO.Out

	fmt.Fprintf(out, "\nTasks in %s:\n\n", cs.Bold(project.Name))

	if len(sections) == 0 {
		fmt.Fprintln(out, "No sections found")
		return nil
	}

	for _, st := range sections {
		fmt.Fprintf(out, "%s:\n", cs.Bold(st.section.Name))

		if len(st.tasks) == 0 {
			fmt.Fprintln(out, "  No tasks in this section")
		} else {
			for i, task := range st.tasks {
				fmt.Fprintf(out, "  %d. %s\n", i+1, task.Name)
			}
		}
		fmt.Fprintln(out)
	}

	return nil
}
