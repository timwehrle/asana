package add

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/prompter"
	taskshared "github.com/timwehrle/asana/pkg/cmd/tasks/shared"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

var followerPropagationDelay = 3 * time.Second

type storyOutput struct {
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

type addOutput struct {
	Story storyOutput `json:"story"`
}

type AddOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Client func() (*asana.Client, error)

	TaskID         string
	Text           string
	MentionUserID  string
	MentionCreator bool
	Output         string
}

func NewCmdAdd(f factory.Factory, runF func(*AddOptions) error) *cobra.Command {
	opts := &AddOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a comment to a task",
		Long:  "Add a new comment to a specific task, optionally mentioning a user or the task creator.",
		Example: heredoc.Doc(`
			$ asana tasks comments add --task 12001234 --text "Can you clarify the acceptance criteria?"
			$ asana tasks comments add --task 12001234 --text "Can you clarify the acceptance criteria?" --mention-creator`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := taskshared.ValidateOutputMode("output", opts.Output); err != nil {
				return err
			}
			if opts.TaskID == "" {
				return fmt.Errorf("--task <gid> is required")
			}
			if strings.TrimSpace(opts.Text) == "" {
				return fmt.Errorf("--text is required")
			}
			if opts.MentionCreator && opts.MentionUserID != "" {
				return fmt.Errorf("--mention-creator and --mention-user are mutually exclusive")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runAdd(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TaskID, "task", "", "Task GID to comment on")
	cmd.Flags().StringVar(&opts.TaskID, "task-id", "", "Task GID to comment on")
	cmd.Flags().StringVar(&opts.Text, "text", "", "Comment text to post")
	cmd.Flags().StringVar(&opts.MentionUserID, "mention-user", "", "Mention a specific user GID in the comment")
	cmd.Flags().BoolVar(&opts.MentionCreator, "mention-creator", false, "Mention the task creator in the comment")
	cmd.Flags().StringVar(&opts.Output, "output", taskshared.OutputText, "Output format: text or json")

	return cmd
}

func runAdd(opts *AddOptions) error {
	client, err := opts.Client()
	if err != nil {
		return err
	}

	task := &asana.Task{ID: opts.TaskID}
	mentionUser, err := resolveMentionUser(client, task, opts)
	if err != nil {
		return err
	}

	story := &asana.StoryBase{
		Text: opts.Text,
	}
	if mentionUser != nil {
		if err := task.AddFollowers(client, []string{mentionUser.ID}); err != nil {
			return fmt.Errorf("failed to add follower %s: %w", mentionUser.ID, err)
		}
		if followerPropagationDelay > 0 {
			time.Sleep(followerPropagationDelay)
		}
		story.Text = ""
		story.HTMLText = buildMentionHTML(mentionUser, opts.Text)
	}

	created, err := task.CreateComment(client, story)
	if err != nil {
		return fmt.Errorf("failed to create comment on task %s: %w", opts.TaskID, err)
	}

	if taskshared.NormalizeOutputMode(opts.Output) == taskshared.OutputJSON {
		return taskshared.WriteJSON(opts.IO.Out, addOutput{Story: toStoryOutput(created)})
	}

	fmt.Fprintf(opts.IO.Out, "%s Comment added to task %s\n", opts.IO.ColorScheme().SuccessIcon, opts.TaskID)
	return nil
}

func resolveMentionUser(client *asana.Client, task *asana.Task, opts *AddOptions) (*asana.User, error) {
	switch {
	case opts.MentionUserID != "":
		return &asana.User{ID: opts.MentionUserID}, nil
	case opts.MentionCreator:
		if err := task.Fetch(client, &asana.Options{
			Fields: []string{"created_by.name", "created_by.gid"},
		}); err != nil {
			return nil, fmt.Errorf("failed to fetch task creator for %s: %w", task.ID, err)
		}
		if task.CreatedBy == nil || task.CreatedBy.ID == "" {
			return nil, fmt.Errorf("task %s does not expose a creator", task.ID)
		}
		return task.CreatedBy, nil
	default:
		return nil, nil
	}
}

func buildMentionHTML(user *asana.User, text string) string {
	name := user.Name
	if name == "" {
		name = user.ID
	}
	return fmt.Sprintf(
		`<body><a data-asana-gid="%s" data-asana-dynamic="false">@%s</a> %s</body>`,
		html.EscapeString(user.ID),
		html.EscapeString(name),
		html.EscapeString(text),
	)
}

func toStoryOutput(story *asana.Story) storyOutput {
	out := storyOutput{
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
