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

	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)
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
			$ asana ts update`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return runUpdate(opts)
		},
	}

	return cmd
}

func runUpdate(opts *UpdateOptions) error {
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
