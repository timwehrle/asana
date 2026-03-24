package attach

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type AttachOptions struct {
	IO *iostreams.IOStreams

	Client func() (*asana.Client, error)

	TaskID string
	URL    string
	Name   string
}

func NewCmdAttach(f factory.Factory, runF func(*AttachOptions) error) *cobra.Command {
	opts := &AttachOptions{
		IO:     f.IOStreams,
		Client: f.Client,
	}

	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach an external URL to a task",
		Long:  "Attach an external resource, such as a pull request URL, to an Asana task.",
		Example: heredoc.Doc(`
			$ asana tasks attach --task 12001234 --url "https://github.com/org/repo/pull/123" --name "PR #123: Short title"`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.TaskID == "" {
				return fmt.Errorf("--task <gid> is required")
			}
			if opts.URL == "" {
				return fmt.Errorf("--url is required")
			}
			if opts.Name == "" {
				return fmt.Errorf("--name is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return runAttach(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TaskID, "task", "", "Task GID to attach the external resource to")
	cmd.Flags().StringVar(&opts.TaskID, "task-id", "", "Task GID to attach the external resource to")
	cmd.Flags().StringVar(&opts.URL, "url", "", "External URL to attach")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Display name for the attachment")

	return cmd
}

func runAttach(opts *AttachOptions) error {
	client, err := opts.Client()
	if err != nil {
		return err
	}

	task := &asana.Task{ID: opts.TaskID}
	attachment, err := task.CreateExternalAttachment(client, &asana.ExternalAttachmentRequest{
		URL:  opts.URL,
		Name: opts.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to attach external URL to task %s: %w", opts.TaskID, err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Attached %s to task %s\n", opts.IO.ColorScheme().SuccessIcon, attachment.Name, opts.TaskID)
	fmt.Fprintf(opts.IO.Out, "  URL: %s\n", opts.URL)
	return nil
}
