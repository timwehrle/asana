package format

import (
	"fmt"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/utils"
	"strings"
	"time"
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

func Tasks(tasks []*asana.Task) []string {
	return formatItems(tasks, func(t *asana.Task) string {
		return fmt.Sprintf("[%s] %s", Date(t.DueOn), t.Name)
	})
}

func Projects(projects []*asana.Project) string {
	names := formatItems(projects, func(p *asana.Project) string {
		return p.Name
	})
	return formatList("Projects: ", names)
}

func Tags(tags []*asana.Tag) string {
	names := formatItems(tags, func(t *asana.Tag) string {
		return t.Name
	})
	return formatList("Tags: ", names)
}

// Notes formats the notes for better readability
func Notes(notes string) string {
	if notes == "" {
		return ""
	}
	return utils.BoldUnderline().Sprintf("Description:") + "\n" + notes + "\n"
}

func Date(date *asana.Date) string {

	if date == nil {
		return "None"
	}

	parsedDate := time.Time(*date)
	location := time.Now().Location()

	parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, location)

	now := time.Now().In(location)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	tomorrow := today.Add(24 * time.Hour)
	weekLater := today.Add(6 * 24 * time.Hour)

	if parsedDate.Equal(today) {
		return "Today"
	}

	if parsedDate.Equal(tomorrow) {
		return "Tomorrow"
	}

	if parsedDate.After(tomorrow) && parsedDate.Before(weekLater) {
		return parsedDate.Format("Mon")
	}

	if parsedDate.Before(today) {
		return parsedDate.Format("Jan 02, 2006")
	}

	return parsedDate.Format("Jan 02, 2006")
}
