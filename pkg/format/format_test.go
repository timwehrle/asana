package format_test

import (
	"testing"
	"time"

	"github.com/timwehrle/asana/pkg/format"

	"github.com/stretchr/testify/assert"
	"github.com/timwehrle/asana/internal/api/asana"
)

type mockStruct struct {
	name string
}

func TestFormatItems(t *testing.T) {
	items := []*mockStruct{
		{name: "item1"},
		{name: "item2"},
	}

	result := format.MapToStrings(items, func(m *mockStruct) string {
		return m.name
	})

	assert.Equal(t, []string{"item1", "item2"}, result)
	assert.Empty(
		t,
		0,
		len(format.MapToStrings([]*mockStruct{}, func(m *mockStruct) string { return m.name })),
	)
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
			prefix:   "MapToStrings: ",
			items:    []string{"one", "two"},
			expected: "MapToStrings: one, two",
		},
		{
			name:     "empty list",
			prefix:   "MapToStrings: ",
			items:    []string{},
			expected: "MapToStrings: None",
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
			result := format.List(tt.prefix, tt.items)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTasks(t *testing.T) {
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

	result := format.Tasks(tasks)
	assert.Len(t, result, 2)
	for _, formattedTask := range result {
		assert.Contains(t, formattedTask, "]")
		assert.Contains(t, formattedTask, "Task")
	}
}

func TestProjects(t *testing.T) {
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

	result := format.Projects(projects)
	assert.Equal(t, "Projects: Project 1, Project 2", result)
	assert.Equal(t, "Projects: None", format.Projects([]*asana.Project{}))
}

func TestTags(t *testing.T) {
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

	result := format.Tags(tags)
	assert.Equal(t, "Tags: Tag 1, Tag 2", result)
	assert.Equal(t, "Tags: None", format.Tags([]*asana.Tag{}))
}

func TestNotes(t *testing.T) {
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
			result := format.Notes(tt.notes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDate(t *testing.T) {
	t.Run("Empty Date", func(t *testing.T) {
		result := format.Date(nil)
		if result != "None" {
			t.Errorf("Expected 'None', got '%s'", result)
		}
	})

	t.Run("Today", func(t *testing.T) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		date := asana.Date(today)
		result := format.Date(&date)
		if result != "Today" {
			t.Errorf("Expected 'Today', got '%s'", result)
		}
	})

	t.Run("Tomorrow", func(t *testing.T) {
		tomorrow := time.Now().Add(24 * time.Hour)
		date := asana.Date(tomorrow)
		result := format.Date(&date)
		if result != "Tomorrow" {
			t.Errorf("Expected 'Tomorrow', got '%s'", result)
		}
	})

	t.Run("Date Within a Week", func(t *testing.T) {
		date := time.Now().Add(3 * 24 * time.Hour)
		expected := date.Format("Mon")
		asanaDate := asana.Date(date)
		result := format.Date(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date After a Week", func(t *testing.T) {
		futureDate := time.Now().Add(8 * 24 * time.Hour)
		expected := futureDate.Format("Jan 02, 2006")
		asanaDate := asana.Date(futureDate)
		result := format.Date(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date Before Today", func(t *testing.T) {
		pastDate := time.Now().Add(8 * (-24) * time.Hour)
		expected := pastDate.Format("Jan 02, 2006")
		asanaDate := asana.Date(pastDate)
		result := format.Date(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}

func TestIndent(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		prefix string
		want   string
	}{
		{
			name:   "non-empty string",
			input:  "line1\nline2\nline3",
			prefix: ">> ",
			want:   ">> line1\n>> line2\n>> line3",
		},
		{
			name:   "empty string",
			input:  "",
			prefix: ">> ",
			want:   "",
		},
		{
			name:   "string with only spaces",
			input:  "   ",
			prefix: ">> ",
			want:   "   ",
		},
		{
			name:   "strings with newlines",
			input:  "line1\n\nline3",
			prefix: "-- ",
			want:   "-- line1\n-- \n-- line3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := format.Indent(tt.input, tt.prefix)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDedent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "string with consistent indentation",
			input: "    line1\n    line2\n    line3",
			want:  "line1\nline2\nline3",
		},
		{
			name:  "string with varying indentation",
			input: "    line1\n  line2\n        line3",
			want:  "  line1\nline2\n      line3",
		},
		{
			name:  "string with no indentation",
			input: "line1\nline2\nline3",
			want:  "line1\nline2\nline3",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "string with blank lines",
			input: "    line1\n\n    line3",
			want:  "line1\n\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := format.Dedent(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
