package create

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	taskshared "github.com/timwehrle/asana/pkg/cmd/tasks/shared"
	"github.com/timwehrle/asana/pkg/convert"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type CreateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter
	Config   func() (*config.Config, error)
	Client   func() (*asana.Client, error)

	Name        string
	Assignee    string
	Due         string
	Description string
	ParentID    string
	ProjectID   string
	SectionID   string
	SectionName string
	DependsOnID string
	Output      string
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

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Task name")
	cmd.Flags().StringVarP(&opts.Assignee, "assignee", "a", "", "Assignee name or 'me'")
	cmd.Flags().StringVarP(&opts.Due, "due", "d", "", "Due date (YYYY-MM-DD, 'today', 'tomorrow')")
	cmd.Flags().StringVarP(&opts.Description, "description", "m", "", "Task description")
	cmd.Flags().StringVar(&opts.ParentID, "parent", "", "Parent task GID to create this task as a subtask")
	cmd.Flags().StringVar(&opts.ProjectID, "project", "", "Project GID to add the task to")
	cmd.Flags().StringVar(&opts.SectionID, "section", "", "Section GID to place the task into")
	cmd.Flags().StringVar(&opts.SectionName, "section-name", "", "Section name to place the task into")
	cmd.Flags().StringVar(&opts.DependsOnID, "depends-on", "", "Task GID that the newly created task should depend on")
	cmd.Flags().StringVar(&opts.Output, "output", taskshared.OutputText, "Output format: text or json")

	return cmd
}

func runCreate(opts *CreateOptions) error {
	if err := taskshared.ValidateOutputMode("output", opts.Output); err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	client, err := opts.Client()
	if err != nil {
		return fmt.Errorf("failed to initialize Asana client: %w", err)
	}

	// Get or prompt for task name
	name := opts.Name
	if name == "" {
		name, err = opts.Prompter.Input("Enter task name: ", "")
		if err != nil {
			return fmt.Errorf("failed to read task name: %w", err)
		}
	}
	if name == "" {
		return fmt.Errorf("task name cannot be empty")
	}

	project, section, err := resolveProjectAndSection(opts, cfg.Workspace.ID, client)
	if err != nil {
		return err
	}

	// Get or prompt for assignee
	assignee, err := getOrSelectAssignee(opts, cfg, client)
	if err != nil {
		return err
	}

	// Get or prompt for due date
	dueDate, err := getOrPromptDueDate(opts)
	if err != nil {
		return err
	}

	description := opts.Description
	if description == "" {
		shouldPromptForDescription, err := opts.Prompter.Confirm("Add description?", "No")
		if err == nil && shouldPromptForDescription {
			description, err = addDescription(opts)
		}
		if err != nil {
			return err
		}
	}

	req := &asana.CreateTaskRequest{
		TaskBase: asana.TaskBase{
			Name:  name,
			DueOn: dueDate,
			Notes: description,
		},
		Workspace: cfg.Workspace.ID,
		Assignee:  assignee.ID,
		Parent:    opts.ParentID,
	}
	if project != nil {
		req.Projects = []string{project.ID}
	}
	if project != nil && section != nil {
		req.Memberships = []*asana.CreateMembership{
			{
				Project: project.ID,
				Section: section.ID,
			},
		}
	}
	if err := req.Validate(); err != nil {
		return fmt.Errorf("task validation failed: %w", err)
	}

	task, err := client.CreateTask(req)
	if err != nil {
		return fmt.Errorf("error creating task: %w", err)
	}

	if opts.DependsOnID != "" {
		if err := task.AddDependencies(client, &asana.AddDependenciesRequest{
			Dependencies: []string{opts.DependsOnID},
		}); err != nil {
			return fmt.Errorf("error adding dependency %s to task %s: %w", opts.DependsOnID, task.ID, err)
		}
	}

	if taskshared.NormalizeOutputMode(opts.Output) == taskshared.OutputJSON {
		return taskshared.WriteJSON(opts.IO.Out, map[string]taskshared.TaskOutput{
			"task": taskshared.ToTaskOutput(task),
		})
	}

	opts.IO.Printf("%s Created task %s\n", cs.SuccessIcon, cs.Bold(task.Name))
	opts.IO.Printf("  %s %s\n", cs.Gray("Assignee:"), assignee.Name)
	if task.DueOn != nil {
		opts.IO.Printf("  %s %s\n", cs.Gray("Due:"), format.Date(task.DueOn))
	}
	if task.PermalinkURL != "" {
		opts.IO.Printf("  %s %s\n", cs.Gray("URL:"), task.PermalinkURL)
	}

	return nil
}

func getOrSelectAssignee(opts *CreateOptions, cfg *config.Config, client *asana.Client) (*asana.User, error) {
	ws := &asana.Workspace{ID: cfg.Workspace.ID}
	users, _, err := ws.Users(client)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch users: %w", err)
	}

	// If flag provided
	if opts.Assignee != "" {
		// Handle 'me' shorthand
		if strings.ToLower(opts.Assignee) == "me" {
			// If no user ID in config, fetch current user
			// This is needed because the user id may not be stored in config yet
			if cfg.UserID == "" {
				currentUser, err := client.CurrentUser()
				if err != nil {
					return nil, fmt.Errorf("failed to fetch current user: %w", err)
				}
				for _, user := range users {
					if user.ID == currentUser.ID {
						return user, nil
					}
				}
				return nil, fmt.Errorf("could not find current user in workspace")
			} else {
				for _, user := range users {
					if user.ID == cfg.UserID {
						return user, nil
					}
				}
				return nil, fmt.Errorf("could not find current user in workspace")
			}
		}

		// Try to match by name
		assigneeLower := strings.ToLower(opts.Assignee)
		for _, user := range users {
			if strings.ToLower(user.Name) == assigneeLower {
				return user, nil
			}
		}

		// Try to match by ID
		for _, user := range users {
			if user.ID == opts.Assignee {
				return user, nil
			}
		}

		return nil, fmt.Errorf("assignee %q not found in workspace", opts.Assignee)
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

func getOrPromptDueDate(opts *CreateOptions) (*asana.Date, error) {
	input := opts.Due
	if input == "" {
		var err error
		input, err = opts.Prompter.Input("Enter due date (YYYY-MM-DD), leave blank for none: ", "")
		if err != nil {
			return nil, fmt.Errorf("failed to read due date: %w", err)
		}
	}
	if input == "" {
		return nil, nil
	}

	now := time.Now()
	switch strings.ToLower(input) {
	case "today":
		return convert.ToDate(now.Format(time.DateOnly), time.DateOnly)
	case "tomorrow":
		return convert.ToDate(now.AddDate(0, 0, 1).Format(time.DateOnly), time.DateOnly)
	}

	due, err := convert.ToDate(input, time.DateOnly)
	if err != nil {
		return nil, fmt.Errorf("invalid due date %q: %w", input, err)
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

func resolveProjectAndSection(
	opts *CreateOptions,
	workspaceID string,
	client *asana.Client,
) (*asana.Project, *asana.Section, error) {
	if opts.SectionID != "" && opts.SectionName != "" {
		return nil, nil, fmt.Errorf("--section and --section-name are mutually exclusive")
	}
	if opts.ProjectID == "" && (opts.SectionID != "" || opts.SectionName != "") {
		return nil, nil, fmt.Errorf("--project is required when using --section or --section-name")
	}

	if opts.ProjectID != "" {
		project := &asana.Project{ID: opts.ProjectID}
		var section *asana.Section
		var err error

		switch {
		case opts.SectionID != "":
			section = &asana.Section{ID: opts.SectionID}
		case opts.SectionName != "":
			section, err = getSectionByName(project.ID, opts.SectionName, client)
			if err != nil {
				return nil, nil, err
			}
		}

		return project, section, nil
	}

	if opts.ParentID != "" {
		return nil, nil, nil
	}

	project, err := getProject(opts, workspaceID, client)
	if err != nil {
		return nil, nil, err
	}
	section, err := getSection(opts, project.ID, client)
	if err != nil {
		return nil, nil, err
	}
	return project, section, nil
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

func getSectionByName(projectID string, sectionName string, client *asana.Client) (*asana.Section, error) {
	project := &asana.Project{ID: projectID}
	sections, _, err := project.Sections(client)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch sections: %w", err)
	}

	for _, section := range sections {
		if strings.EqualFold(strings.TrimSpace(section.Name), strings.TrimSpace(sectionName)) {
			return section, nil
		}
	}

	return nil, fmt.Errorf("section %q not found in project %s", sectionName, projectID)
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
