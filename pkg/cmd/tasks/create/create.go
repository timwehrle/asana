package create

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/convert"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
	"strings"
	"time"
)

type CreateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter
	Config   func() (*config.Config, error)
	Client   func() (*asana.Client, error)
}

func NewCmdCreate(f factory.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new task",
		Long:  "Create a new task in Asana.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runCreate(opts)
		},
	}

	return cmd
}

func runCreate(opts *CreateOptions) error {
	cs := opts.IO.ColorScheme()

	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	client, err := opts.Client()
	if err != nil {
		return fmt.Errorf("failed to initialize Asana client: %w", err)
	}

	// Prompt for task name
	name, err := opts.Prompter.Input("Enter task name: ", "")
	if err != nil {
		return fmt.Errorf("failed to read task name: %w", err)
	}
	if name == "" {
		return fmt.Errorf("task name cannot be empty")
	}

	// Prompt for assignee
	assignee, err := selectAssignee(opts, cfg.Workspace.ID, client)
	if err != nil {
		return err
	}

	// Prompt for due date
	dueDate, err := getDueDate(opts)
	if err != nil {
		return err
	}

	// Prompt for task description
	description, err := addDescription(opts)
	if err != nil {
		return err
	}

	// Prompt for project
	project, err := getProject(opts, cfg.Workspace.ID, client)
	if err != nil {
		return err
	}

	// Prompt for section
	section, err := getSection(opts, project.ID, client)
	if err != nil {
		return err
	}

	req := &asana.CreateTaskRequest{
		TaskBase: asana.TaskBase{
			Name:  name,
			DueOn: dueDate,
			Notes: description,
		},
		Workspace: cfg.Workspace.ID,
		Assignee:  assignee.ID,

		// Currently only one project is supported
		Projects: []string{project.ID},

		// Both project and section ID are expected
		Memberships: []*asana.CreateMembership{
			{
				Project: project.ID,
				Section: section.ID,
			},
		},
	}
	if err := req.Validate(); err != nil {
		return fmt.Errorf("task validation failed: %w", err)
	}

	task, err := client.CreateTask(req)
	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	opts.IO.Printf("%s Created task %q with due date %s\n", cs.SuccessIcon, task.Name, strings.ToLower(format.Date(task.DueOn)))
	return nil
}

func selectAssignee(opts *CreateOptions, workspaceID string, client *asana.Client) (*asana.User, error) {
	ws := &asana.Workspace{ID: workspaceID}
	users, _, err := ws.Users(client)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch users: %w", err)
	}

	names := format.MapToStrings(users, func(u *asana.User) string {
		return u.Name
	})

	selected, err := opts.Prompter.Select("Select assignee: ", names)
	if err != nil {
		return nil, fmt.Errorf("assignee selection failed: %w", err)
	}
	return users[selected], nil
}

func getDueDate(opts *CreateOptions) (*asana.Date, error) {
	input, err := opts.Prompter.Input("Enter due date (YYYY-MM-DD), leave blank for none: ", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read due date: %w", err)
	}

	var due *asana.Date
	if input != "" {
		due, err = convert.ToDate(input, time.DateOnly)
		if err != nil {
			return nil, fmt.Errorf("invalid due date %q: %w", input, err)
		}
	}
	return due, nil
}

func addDescription(opts *CreateOptions) (string, error) {
	description, err := opts.Prompter.Editor("Enter task description: ", "")
	if err != nil {
		return "", fmt.Errorf("failed to read task description: %w", err)
	}

	return strings.TrimSpace(description), nil
}

func getProject(opts *CreateOptions, workspaceID string, client *asana.Client) (*asana.Project, error) {
	ws := &asana.Workspace{ID: workspaceID}
	projects, err := ws.AllProjects(client)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch projects: %w", err)
	}

	names := format.MapToStrings(projects, func(p *asana.Project) string {
		return p.Name
	})

	selected, err := opts.Prompter.Select("Select project: ", names)
	if err != nil {
		return nil, fmt.Errorf("project selection failed: %w", err)
	}
	return projects[selected], nil
}

func getSection(opts *CreateOptions, projectID string, client *asana.Client) (*asana.Section, error) {
	project := &asana.Project{ID: projectID}
	sections, _, err := project.Sections(client)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch sections: %w", err)
	}

	names := format.MapToStrings(sections, func(p *asana.Section) string {
		return p.Name
	})

	selected, err := opts.Prompter.Select("Select section: ", names)
	if err != nil {
		return nil, fmt.Errorf("section selection failed: %w", err)
	}
	return sections[selected], nil
}
