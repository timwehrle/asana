package update

import (
	"fmt"
	"strings"
	"time"

	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/timwehrle/asana/internal/api/asana"
	taskshared "github.com/timwehrle/asana/pkg/cmd/tasks/shared"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/convert"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/format"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type UpdateAction int

const (
	ActionComplete UpdateAction = iota
	ActionEditName
	ActionEditDescription
	ActionSetDueDate
	ActionCancel
)

type taskAction struct {
	name   string
	action UpdateAction
}

var availableActions = []taskAction{
	{name: "Mark as Completed", action: ActionComplete},
	{name: "Edit Task Name", action: ActionEditName},
	{name: "Edit Description", action: ActionEditDescription},
	{name: "Set Due Date", action: ActionSetDueDate},
	{name: "Cancel", action: ActionCancel},
}

type UpdateOptions struct {
	IO       *iostreams.IOStreams
	Prompter prompter.Prompter

	Config           func() (*config.Config, error)
	Client           func() (*asana.Client, error)
	TaskID           string
	NewName          string
	Notes            string
	NotesFile        string
	PrependNotes     string
	PrependNotesFile string
	Complete         bool
	DueOn            string
	Output           string
}

func NewCmdUpdate(f factory.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:       f.IOStreams,
		Prompter: f.Prompter,
		Config:   f.Config,
		Client:   f.Client,
	}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update details of a specific task",
		Long:  "Retrieve task details and select one for updating it.",
		Args:  cobra.NoArgs,
		Example: heredoc.Doc(`
			$ asana tasks update
			$ asana ts update
			$ asana tasks update --task 12001234 --prepend-notes "Branch: codex/my-branch\n\n" --output json`),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := taskshared.ValidateOutputMode("output", opts.Output); err != nil {
				return err
			}
			if err := validateDirectUpdateOptions(opts); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runUpdate(opts)
		},
	}
	cmd.Flags().StringVar(&opts.TaskID, "task", "", "Task GID to update directly without prompting")
	cmd.Flags().StringVar(&opts.TaskID, "task-id", "", "Task GID to update directly without prompting")
	cmd.Flags().StringVar(&opts.NewName, "name", "", "New task name")
	cmd.Flags().StringVar(&opts.Notes, "notes", "", "Replace the task description")
	cmd.Flags().StringVar(&opts.NotesFile, "notes-file", "", "Replace the task description with file contents")
	cmd.Flags().StringVar(&opts.PrependNotes, "prepend-notes", "", "Prepend text to the existing task description")
	cmd.Flags().StringVar(&opts.PrependNotesFile, "prepend-notes-file", "", "Prepend file contents to the existing task description")
	cmd.Flags().BoolVar(&opts.Complete, "complete", false, "Mark the task as completed")
	cmd.Flags().StringVar(&opts.DueOn, "due-on", "", "Set the task due date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.Output, "output", taskshared.OutputText, "Output format: text or json")

	return cmd
}

func runUpdate(opts *UpdateOptions) error {
	if err := validateDirectUpdateOptions(opts); err != nil {
		return err
	}

	if hasDirectUpdateFlags(opts) {
		if opts.TaskID == "" {
			return fmt.Errorf("non-interactive updates require --task <gid>")
		}
		return runDirectUpdate(opts)
	}

	if opts.TaskID != "" {
		if err := taskshared.EnsureInteractiveAllowed(opts.IO, "an action flag such as --notes, --prepend-notes, --name, --complete, or --due-on"); err != nil {
			return err
		}

		task, err := fetchTaskForUpdate(opts)
		if err != nil {
			return err
		}
		action, err := selectAction(opts)
		if err != nil {
			return err
		}
		if err := performAction(opts, task, action); err != nil {
			return fmt.Errorf("failed to perform action: %w", err)
		}
		return nil
	}

	if err := taskshared.EnsureInteractiveAllowed(opts.IO, "--task <gid> with an action flag"); err != nil {
		return err
	}

	task, err := selectTask(opts)
	if err != nil {
		return err
	}

	action, err := selectAction(opts)
	if err != nil {
		return err
	}

	if err := performAction(opts, task, action); err != nil {
		return fmt.Errorf("failed to perform action: %w", err)
	}

	return nil
}

func validateDirectUpdateOptions(opts *UpdateOptions) error {
	if err := cmdutils.ValidateDate("due-on", opts.DueOn); err != nil {
		return err
	}
	if opts.Notes != "" && opts.NotesFile != "" {
		return fmt.Errorf("--notes and --notes-file are mutually exclusive")
	}
	if opts.PrependNotes != "" && opts.PrependNotesFile != "" {
		return fmt.Errorf("--prepend-notes and --prepend-notes-file are mutually exclusive")
	}
	if (opts.Notes != "" || opts.NotesFile != "") && (opts.PrependNotes != "" || opts.PrependNotesFile != "") {
		return fmt.Errorf("replace and prepend note operations are mutually exclusive")
	}
	return nil
}

func hasDirectUpdateFlags(opts *UpdateOptions) bool {
	return opts.NewName != "" ||
		opts.Notes != "" ||
		opts.NotesFile != "" ||
		opts.PrependNotes != "" ||
		opts.PrependNotesFile != "" ||
		opts.Complete ||
		opts.DueOn != ""
}

func selectTask(opts *UpdateOptions) (*asana.Task, error) {
	cfg, err := opts.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	client, err := opts.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create Asana client: %w", err)
	}

	tasks, _, err := client.QueryTasks(&asana.TaskQuery{
		Assignee:       "me",
		Workspace:      cfg.Workspace.ID,
		CompletedSince: "now",
	}, &asana.Options{
		Fields: []string{"name", "due_on"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	if len(tasks) == 0 {
		fmt.Fprintln(opts.IO.Out, "No tasks found.")
		return nil, nil
	}

	taskNames := format.Tasks(tasks)
	index, err := opts.Prompter.Select("Select the task to update:", taskNames)
	if err != nil {
		return nil, fmt.Errorf("failed to select task: %w", err)
	}

	selectedTask := tasks[index]
	if err := selectedTask.Fetch(client); err != nil {
		return nil, fmt.Errorf("failed to fetch task details: %w", err)
	}

	return selectedTask, nil
}

func fetchTaskForUpdate(opts *UpdateOptions) (*asana.Task, error) {
	client, err := opts.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create Asana client: %w", err)
	}

	return taskshared.FetchTaskByID(client, opts.TaskID, []string{
		"name",
		"notes",
		"completed",
		"due_on",
		"permalink_url",
		"projects.name",
		"tags.name",
	})
}

func selectAction(opts *UpdateOptions) (UpdateAction, error) {
	actions := make([]string, len(availableActions))
	for i, action := range availableActions {
		actions[i] = action.name
	}

	index, err := opts.Prompter.Select("What do you want to do with this task:", actions)
	if err != nil {
		return 0, fmt.Errorf("failed to select action: %w", err)
	}

	return availableActions[index].action, nil
}

func performAction(opts *UpdateOptions, task *asana.Task, action UpdateAction) error {
	client, err := opts.Client()
	if err != nil {
		return fmt.Errorf("failed to create Asana client: %w", err)
	}

	cs := opts.IO.ColorScheme()

	switch action {
	case ActionComplete:
		return completeTask(client, task, cs)
	case ActionEditName:
		return editTaskName(opts, client, task, cs)
	case ActionEditDescription:
		return editTaskDescription(opts, client, task, cs)
	case ActionSetDueDate:
		return setTaskDueDate(opts, client, task, cs)
	case ActionCancel:
		fmt.Fprintf(
			opts.IO.Out,
			"%s Operation canceled. You can rerun the command to try again.\n",
			cs.SuccessIcon,
		)
		return nil
	default:
		return fmt.Errorf("unknown action: %d", action)
	}
}

func runDirectUpdate(opts *UpdateOptions) error {
	task, err := fetchTaskForUpdate(opts)
	if err != nil {
		return err
	}

	client, err := opts.Client()
	if err != nil {
		return fmt.Errorf("failed to create Asana client: %w", err)
	}

	request, updatedFields, err := buildDirectUpdateRequest(opts, task)
	if err != nil {
		return err
	}
	if len(updatedFields) == 0 {
		return fmt.Errorf("non-interactive updates require at least one action flag")
	}

	if err := task.Update(client, request); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if taskshared.NormalizeOutputMode(opts.Output) == taskshared.OutputJSON {
		return taskshared.WriteJSON(opts.IO.Out, taskshared.TaskUpdateOutput{
			Task:          taskshared.ToTaskOutput(task),
			UpdatedFields: updatedFields,
		})
	}

	fmt.Fprintf(opts.IO.Out, "%s Task updated\n", opts.IO.ColorScheme().SuccessIcon)
	return nil
}

func buildDirectUpdateRequest(opts *UpdateOptions, task *asana.Task) (*asana.UpdateTaskRequest, []string, error) {
	request := &asana.UpdateTaskRequest{}
	updatedFields := make([]string, 0, 4)

	if opts.NewName != "" {
		request.Name = opts.NewName
		updatedFields = append(updatedFields, "name")
	}

	notes, notesField, err := resolveNotesUpdate(opts, task.Notes)
	if err != nil {
		return nil, nil, err
	}
	if notesField != "" {
		request.Notes = notes
		updatedFields = append(updatedFields, notesField)
	}

	if opts.Complete {
		request.Completed = asana.Bool(true)
		updatedFields = append(updatedFields, "completed")
	}

	if opts.DueOn != "" {
		dueOn, err := convert.ToDate(opts.DueOn, time.DateOnly)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid date format: %w", err)
		}
		request.DueOn = dueOn
		updatedFields = append(updatedFields, "due_on")
	}

	return request, updatedFields, nil
}

func resolveNotesUpdate(opts *UpdateOptions, existingNotes string) (string, string, error) {
	switch {
	case opts.Notes != "":
		return opts.Notes, "notes", nil
	case opts.NotesFile != "":
		content, err := taskshared.ReadFile(opts.NotesFile)
		if err != nil {
			return "", "", err
		}
		return content, "notes", nil
	case opts.PrependNotes != "":
		return taskshared.PrependNotes(opts.PrependNotes, existingNotes), "notes", nil
	case opts.PrependNotesFile != "":
		content, err := taskshared.ReadFile(opts.PrependNotesFile)
		if err != nil {
			return "", "", err
		}
		return taskshared.PrependNotes(content, existingNotes), "notes", nil
	default:
		return "", "", nil
	}
}

func completeTask(client *asana.Client, task *asana.Task, cs *iostreams.ColorScheme) error {
	completed := true
	updateRequest := &asana.UpdateTaskRequest{
		TaskBase: asana.TaskBase{
			Completed: &completed,
		},
	}

	if err := task.Update(client, updateRequest); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	fmt.Printf("%s Task completed\n", cs.SuccessIcon)

	return nil
}

func editTaskName(
	opts *UpdateOptions,
	client *asana.Client,
	task *asana.Task,
	cs *iostreams.ColorScheme,
) error {
	newName, err := opts.Prompter.Input("Enter the new task name:", task.Name)
	if err != nil {
		return fmt.Errorf("failed to get input: %w", err)
	}

	newName = strings.TrimSpace(newName)
	if newName == task.Name {
		fmt.Fprintf(opts.IO.Out, "%s No changes made to task name\n", cs.WarningIcon)
		return nil
	}

	updateRequest := &asana.UpdateTaskRequest{
		TaskBase: asana.TaskBase{
			Name: newName,
		},
	}

	if err := task.Update(client, updateRequest); err != nil {
		return fmt.Errorf("failed to update task name: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Task name updated\n", cs.SuccessIcon)
	return nil
}

func editTaskDescription(
	opts *UpdateOptions,
	client *asana.Client,
	task *asana.Task,
	cs *iostreams.ColorScheme,
) error {
	existingDescription := strings.TrimSpace(task.Notes)
	newDescription, err := opts.Prompter.Editor("Edit the description:", existingDescription)
	if err != nil {
		return fmt.Errorf("failed to get input: %w", err)
	}

	newDescription = strings.TrimSpace(newDescription)
	if newDescription == existingDescription {
		fmt.Fprintf(opts.IO.Out, "%s No changes made to description\n", cs.WarningIcon)
		return nil
	}

	updateRequest := &asana.UpdateTaskRequest{
		TaskBase: asana.TaskBase{
			Notes: newDescription,
		},
	}

	if err = task.Update(client, updateRequest); err != nil {
		return fmt.Errorf("failed to update task description: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Description updated\n", cs.SuccessIcon)
	return nil
}

func setTaskDueDate(
	opts *UpdateOptions,
	client *asana.Client,
	task *asana.Task,
	cs *iostreams.ColorScheme,
) error {
	input, err := opts.Prompter.Input(
		"Enter the new due date (YYYY-MM-DD):",
		format.Date(task.DueOn),
	)
	if err != nil {
		return fmt.Errorf("failed to get input: %w", err)
	}

	dueDate, err := convert.ToDate(input, time.DateOnly)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	updateRequest := &asana.UpdateTaskRequest{
		TaskBase: asana.TaskBase{
			DueOn: dueDate,
		},
	}

	if err := task.Update(client, updateRequest); err != nil {
		return fmt.Errorf("failed to update task due date: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Due date updated\n", cs.SuccessIcon)
	return nil
}
