package shared

import (
	"github.com/stretchr/testify/assert"
	"github.com/timwehrle/asana-go"
	"testing"
	"time"
)

type mockStruct struct {
	name string
}

func TestFormatItems(t *testing.T) {
	items := []*mockStruct{
		{name: "item1"},
		{name: "item2"},
	}

	result := formatItems(items, func(m *mockStruct) string {
		return m.name
	})

	assert.Equal(t, []string{"item1", "item2"}, result)
	assert.Equal(t, 0, len(formatItems([]*mockStruct{}, func(m *mockStruct) string { return m.name })))
}

func TestFormatList(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		items    []string
		expected string
	}{
		{
			name:     "normal list",
			prefix:   "Items: ",
			items:    []string{"one", "two"},
			expected: "Items: one, two",
		},
		{
			name:     "empty list",
			prefix:   "Items: ",
			items:    []string{},
			expected: "Items: None",
		},
		{
			name:     "single item",
			prefix:   "Item: ",
			items:    []string{"one"},
			expected: "Item: one",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatList(tt.prefix, tt.items)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatTasks(t *testing.T) {
	now := time.Now()
	dueDate := asana.Date(now)
	tasks := []*asana.Task{
		{
			TaskBase: asana.TaskBase{
				Name:  "Task 1",
				DueOn: &dueDate,
			},
		},
		{
			TaskBase: asana.TaskBase{
				Name:  "Task 2",
				DueOn: &dueDate,
			},
		},
	}

	result := FormatTasks(tasks)
	assert.Len(t, result, 2)
	for _, formattedTask := range result {
		assert.Contains(t, formattedTask, "]")
		assert.Contains(t, formattedTask, "Task")
	}
}

func TestFormatProjects(t *testing.T) {
	projects := []*asana.Project{
		{
			ProjectBase: asana.ProjectBase{
				Name: "Project 1",
			},
		},
		{
			ProjectBase: asana.ProjectBase{
				Name: "Project 2",
			},
		},
	}

	result := FormatProjects(projects)
	assert.Equal(t, "Projects: Project 1, Project 2", result)
	assert.Equal(t, "Projects: None", FormatProjects([]*asana.Project{}))
}

func TestFormatTags(t *testing.T) {
	tags := []*asana.Tag{
		{
			TagBase: asana.TagBase{
				Name: "Tag 1",
			},
		},
		{
			TagBase: asana.TagBase{
				Name: "Tag 2",
			},
		},
	}

	result := FormatTags(tags)
	assert.Equal(t, "Tags: Tag 1, Tag 2", result)
	assert.Equal(t, "Tags: None", FormatTags([]*asana.Tag{}))
}

func TestFormatNotes(t *testing.T) {
	tests := []struct {
		name     string
		notes    string
		expected string
	}{
		{
			name:     "with content",
			notes:    "Some notes",
			expected: "Description:\nSome notes\n",
		},
		{
			name:     "empty notes",
			notes:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNotes(tt.notes)
			assert.Equal(t, tt.expected, result)
		})
	}
}
