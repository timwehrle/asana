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
	GID          string                  `json:"gid"`
	Name         string                  `json:"name,omitempty"`
	Notes        string                  `json:"notes,omitempty"`
	Completed    bool                    `json:"completed"`
	DueOn        string                  `json:"due_on,omitempty"`
	PermalinkURL string                  `json:"permalink_url,omitempty"`
	Dependencies []TaskDependencyOutput  `json:"dependencies,omitempty"`
	CustomFields []TaskCustomFieldOutput `json:"custom_fields,omitempty"`
	Memberships  []TaskMembershipOutput  `json:"memberships,omitempty"`
}

type TaskMembershipOutput struct {
	Project *TaskResourceOutput `json:"project,omitempty"`
	Section *TaskResourceOutput `json:"section,omitempty"`
}

type TaskResourceOutput struct {
	GID  string `json:"gid"`
	Name string `json:"name,omitempty"`
}

type TaskDependencyOutput struct {
	GID       string `json:"gid"`
	Name      string `json:"name,omitempty"`
	Completed bool   `json:"completed"`
}

type TaskCustomFieldOutput struct {
	GID             string                `json:"gid"`
	Name            string                `json:"name,omitempty"`
	Type            string                `json:"type,omitempty"`
	TextValue       *string               `json:"text_value"`
	NumberValue     *float64              `json:"number_value"`
	BooleanValue    *bool                 `json:"boolean_value"`
	DisplayValue    *string               `json:"display_value"`
	EnumValue       *TaskEnumValueOutput  `json:"enum_value"`
	MultiEnumValues []TaskEnumValueOutput `json:"multi_enum_values,omitempty"`
}

type TaskEnumValueOutput struct {
	GID   string `json:"gid"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
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
	if len(task.Dependencies) > 0 {
		out.Dependencies = make([]TaskDependencyOutput, 0, len(task.Dependencies))
		for _, dependency := range task.Dependencies {
			if dependency == nil {
				continue
			}

			dependencyOut := TaskDependencyOutput{
				GID:  dependency.ID,
				Name: dependency.Name,
			}
			if dependency.Completed != nil {
				dependencyOut.Completed = *dependency.Completed
			}

			out.Dependencies = append(out.Dependencies, dependencyOut)
		}
	}
	if len(task.CustomFields) > 0 {
		out.CustomFields = make([]TaskCustomFieldOutput, 0, len(task.CustomFields))
		for _, field := range task.CustomFields {
			if field == nil {
				continue
			}

			fieldOut := TaskCustomFieldOutput{
				GID:          field.ID,
				Name:         field.Name,
				Type:         string(field.ResourceSubtype),
				TextValue:    field.TextValue,
				NumberValue:  field.NumberValue,
				BooleanValue: field.BooleanValue,
				DisplayValue: field.DisplayValue,
			}
			if field.EnumValue != nil {
				fieldOut.EnumValue = &TaskEnumValueOutput{
					GID:   field.EnumValue.ID,
					Name:  field.EnumValue.Name,
					Color: field.EnumValue.Color,
				}
			}
			if len(field.MultiEnumValues) > 0 {
				fieldOut.MultiEnumValues = make([]TaskEnumValueOutput, 0, len(field.MultiEnumValues))
				for _, value := range field.MultiEnumValues {
					if value == nil {
						continue
					}
					fieldOut.MultiEnumValues = append(fieldOut.MultiEnumValues, TaskEnumValueOutput{
						GID:   value.ID,
						Name:  value.Name,
						Color: value.Color,
					})
				}
			}

			out.CustomFields = append(out.CustomFields, fieldOut)
		}
	}
	if len(task.Memberships) > 0 {
		out.Memberships = make([]TaskMembershipOutput, 0, len(task.Memberships))
		for _, membership := range task.Memberships {
			if membership == nil {
				continue
			}

			membershipOut := TaskMembershipOutput{}
			if membership.Project != nil {
				membershipOut.Project = &TaskResourceOutput{
					GID:  membership.Project.ID,
					Name: membership.Project.Name,
				}
			}
			if membership.Section != nil {
				membershipOut.Section = &TaskResourceOutput{
					GID:  membership.Section.ID,
					Name: membership.Section.Name,
				}
			}

			out.Memberships = append(out.Memberships, membershipOut)
		}
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
