package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"github.com/timwehrle/asana/pkg/iostreams"
)

const (
	OutputText = "text"
	OutputJSON = "json"
)

var validOutputModes = []string{OutputJSON, OutputText}

type TaskOutput struct {
	GID          string `json:"gid"`
	Name         string `json:"name,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Completed    bool   `json:"completed"`
	DueOn        string `json:"due_on,omitempty"`
	PermalinkURL string `json:"permalink_url,omitempty"`
}

type TaskListOutput struct {
	Tasks          []TaskOutput `json:"tasks"`
	NextPageOffset string       `json:"next_page_offset,omitempty"`
}

type TaskUpdateOutput struct {
	Task          TaskOutput `json:"task"`
	UpdatedFields []string   `json:"updated_fields,omitempty"`
}

func ValidateOutputMode(flagName, value string) error {
	if value == "" {
		return nil
	}
	return cmdutils.ValidateStringEnum(flagName, value, validOutputModes)
}

func NormalizeOutputMode(value string) string {
	if value == "" {
		return OutputText
	}
	return value
}

func EnsureInteractiveAllowed(io *iostreams.IOStreams, requiredFlags string) error {
	if io != nil && !io.IsStdinTTY {
		return fmt.Errorf("interactive prompts are unavailable in non-interactive mode; provide %s", requiredFlags)
	}
	return nil
}

func FetchTaskByID(client *asana.Client, taskID string, fields []string) (*asana.Task, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	task := &asana.Task{ID: taskID}
	if err := task.Fetch(client, &asana.Options{Fields: fields}); err != nil {
		return nil, fmt.Errorf("failed to fetch task %s: %w", taskID, err)
	}

	return task, nil
}

func ToTaskOutput(task *asana.Task) TaskOutput {
	out := TaskOutput{
		GID:          task.ID,
		Name:         task.Name,
		Notes:        task.Notes,
		PermalinkURL: task.PermalinkURL,
	}

	if task.Completed != nil {
		out.Completed = *task.Completed
	}
	if task.DueOn != nil {
		out.DueOn = time.Time(*task.DueOn).Format(time.DateOnly)
	}

	return out
}

func WriteJSON(w io.Writer, payload any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(payload)
}

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}
	return string(content), nil
}

func PrependNotes(prefix, existing string) string {
	if existing == "" {
		return prefix
	}
	return prefix + existing
}
