package move

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type MoveOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Client func() (*asana.Client, error)

	TaskID      string
	ProjectID   string
	SectionID   string
	SectionName string
}

func NewCmdMove(f factory.Factory, runF func(*MoveOptions) error) *cobra.Command {
	opts := &MoveOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "move",
		Short: "Move a task into a project or section",
		Long:  "Add a task to a project and optionally place it into a specific section for triage.",
		Example: heredoc.Doc(`
			$ asana tasks move --task 12001234 --project 12009999 --section-name "Backlog"
			$ asana tasks move --task 12001234 --project 12009999 --section 12008888`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.TaskID == "" {
				return fmt.Errorf("--task <gid> is required")
			}
			if opts.ProjectID == "" {
				return fmt.Errorf("--project <gid> is required")
			}
			if opts.SectionID != "" && opts.SectionName != "" {
				return fmt.Errorf("--section and --section-name are mutually exclusive")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runMove(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TaskID, "task", "", "Task GID to move")
	cmd.Flags().StringVar(&opts.TaskID, "task-id", "", "Task GID to move")
	cmd.Flags().StringVar(&opts.ProjectID, "project", "", "Project GID to move the task into")
	cmd.Flags().StringVar(&opts.SectionID, "section", "", "Section GID within the project")
	cmd.Flags().StringVar(&opts.SectionName, "section-name", "", "Section name within the project")

	return cmd
}

func runMove(opts *MoveOptions) error {
	client, err := opts.Client()
	if err != nil {
		return err
	}

	sectionID := opts.SectionID
	if sectionID == "" && opts.SectionName != "" {
		sectionID, err = resolveSectionID(client, opts.ProjectID, opts.SectionName)
		if err != nil {
			return err
		}
	}

	task := &asana.Task{ID: opts.TaskID}
	if err := task.AddProject(client, &asana.AddProjectRequest{
		Project: opts.ProjectID,
		Section: sectionID,
	}); err != nil {
		return fmt.Errorf("failed to move task %s: %w", opts.TaskID, err)
	}

	if sectionID != "" {
		fmt.Fprintf(opts.IO.Out, "%s Task moved to project %s / section %s\n", opts.IO.ColorScheme().SuccessIcon, opts.ProjectID, sectionID)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Task moved to project %s\n", opts.IO.ColorScheme().SuccessIcon, opts.ProjectID)
	}
	return nil
}

func resolveSectionID(client *asana.Client, projectID, sectionName string) (string, error) {
	project := &asana.Project{ID: projectID}
	sections, _, err := project.Sections(client, &asana.Options{Fields: []string{"name"}})
	if err != nil {
		return "", fmt.Errorf("failed to fetch sections for project %s: %w", projectID, err)
	}

	for _, section := range sections {
		if strings.EqualFold(strings.TrimSpace(section.Name), strings.TrimSpace(sectionName)) {
			return section.ID, nil
		}
	}
	return "", fmt.Errorf("section %q not found in project %s", sectionName, projectID)
}
