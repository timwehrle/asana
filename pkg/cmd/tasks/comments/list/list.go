package list

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/prompter"
	taskshared "github.com/timwehrle/asana/pkg/cmd/tasks/shared"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type commentOutput struct {
	ID        string      `json:"id"`
	Text      string      `json:"text"`
	HTMLText  string      `json:"html_text,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	CreatedBy *userOutput `json:"created_by,omitempty"`
}

type userOutput struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type listOutput struct {
	TaskID   string          `json:"task_id"`
	Comments []commentOutput `json:"comments"`
}

type ListOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Client func() (*asana.Client, error)

	TaskID string
	Output string
}

func NewCmdList(f factory.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List comments on a task",
		Long: heredoc.Doc(`
			List comment stories on a specific task.
		`),
		Example: heredoc.Doc(`
			$ asana tasks comments list --task 12001234
			$ asana tasks comments list --task 12001234 --output json`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := taskshared.ValidateOutputMode("output", opts.Output); err != nil {
				return err
			}
			if opts.TaskID == "" {
				return fmt.Errorf("--task <gid> is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runList(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TaskID, "task", "", "Task GID to read comments from")
	cmd.Flags().StringVar(&opts.TaskID, "task-id", "", "Task GID to read comments from")
	cmd.Flags().StringVar(&opts.Output, "output", taskshared.OutputText, "Output format: text or json")

	return cmd
}

func runList(opts *ListOptions) error {
	client, err := opts.Client()
	if err != nil {
		return err
	}

	task := &asana.Task{ID: opts.TaskID}
	stories, _, err := task.Stories(client, &asana.Options{
		Fields: []string{
			"created_by.name",
			"created_by.gid",
			"created_at",
			"text",
			"html_text",
			"resource_subtype",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to fetch comments for task %s: %w", opts.TaskID, err)
	}

	comments := filterCommentStories(stories)
	if taskshared.NormalizeOutputMode(opts.Output) == taskshared.OutputJSON {
		payload := listOutput{
			TaskID:   opts.TaskID,
			Comments: make([]commentOutput, 0, len(comments)),
		}
		for _, story := range comments {
			payload.Comments = append(payload.Comments, toCommentOutput(story))
		}
		return taskshared.WriteJSON(opts.IO.Out, payload)
	}

	fmt.Fprintf(opts.IO.Out, "\nComments on task %s:\n\n", opts.TaskID)
	if len(comments) == 0 {
		fmt.Fprintln(opts.IO.Out, "No comments found")
		return nil
	}
	for i, story := range comments {
		author := "Unknown"
		if story.CreatedBy != nil && story.CreatedBy.Name != "" {
			author = story.CreatedBy.Name
		}
		createdAt := ""
		if story.CreatedAt != nil {
			createdAt = story.CreatedAt.In(time.Local).Format("Jan 02, 2006 15:04")
		}
		if createdAt != "" {
			fmt.Fprintf(opts.IO.Out, "%d. [%s] %s: %s\n", i+1, createdAt, author, story.Text)
		} else {
			fmt.Fprintf(opts.IO.Out, "%d. %s: %s\n", i+1, author, story.Text)
		}
	}
	return nil
}

func filterCommentStories(stories []*asana.Story) []*asana.Story {
	comments := make([]*asana.Story, 0, len(stories))
	for _, story := range stories {
		if story.ResourceSubtype == "comment_added" {
			comments = append(comments, story)
		}
	}
	return comments
}

func toCommentOutput(story *asana.Story) commentOutput {
	out := commentOutput{
		ID:       story.ID,
		Text:     story.Text,
		HTMLText: story.HTMLText,
	}
	if story.CreatedAt != nil {
		out.CreatedAt = story.CreatedAt.Format(time.RFC3339)
	}
	if story.CreatedBy != nil {
		out.CreatedBy = &userOutput{
			GID:  story.CreatedBy.ID,
			Name: story.CreatedBy.Name,
		}
	}
	return out
}
