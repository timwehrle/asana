package shared

import (
	"fmt"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/utils"
	"strings"
)

func formatItems[T any](items []*T, nameFunc func(*T) string) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = nameFunc(item)
	}
	return names
}

func formatList(prefix string, items []string) string {
	if len(items) == 0 {
		return prefix + "None"
	}
	return prefix + strings.Join(items, ", ")
}

func FormatTasks(tasks []*asana.Task) []string {
	return formatItems(tasks, func(t *asana.Task) string {
		return fmt.Sprintf("[%s] %s", utils.FormatDate(t.DueOn), t.Name)
	})
}

func FormatProjects(projects []*asana.Project) string {
	names := formatItems(projects, func(p *asana.Project) string {
		return p.Name
	})
	return formatList("Projects: ", names)
}

func FormatTags(tags []*asana.Tag) string {
	names := formatItems(tags, func(t *asana.Tag) string {
		return t.Name
	})
	return formatList("Tags: ", names)
}

// FormatNotes formats the notes for better readability
func FormatNotes(notes string) string {
	if notes == "" {
		return ""
	}
	return utils.BoldUnderline().Sprintf("Description:") + "\n" + notes + "\n"
}
