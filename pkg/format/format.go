package format

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/pkg/iostreams"
)

func MapToStrings[T any](items []*T, fn func(*T) string) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = fn(item)
	}
	return out
}

func List(prefix string, items []string) string {
	if len(items) == 0 {
		return prefix + "None"
	}
	return prefix + strings.Join(items, ", ")
}

func Tasks(tasks []*asana.Task) []string {
	return MapToStrings(tasks, func(t *asana.Task) string {
		return fmt.Sprintf("[%s] %s", Date(t.DueOn), t.Name)
	})
}

func Projects(projects []*asana.Project) string {
	names := MapToStrings(projects, func(p *asana.Project) string {
		return p.Name
	})
	return List("Projects: ", names)
}

func Tags(tags []*asana.Tag) string {
	names := MapToStrings(tags, func(t *asana.Tag) string {
		return t.Name
	})
	return List("Tags: ", names)
}

// Notes formats the notes for better readability
func Notes(notes string) string {
	io := iostreams.System()
	cs := io.ColorScheme()

	if notes == "" {
		return ""
	}
	return cs.Bold("Description:") + "\n" + notes + "\n"
}

func Date(date *asana.Date) string {
	if date == nil {
		return "None"
	}

	parsedDate := time.Time(*date)
	location := time.Now().Location()

	parsedDate = time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		0,
		0,
		0,
		0,
		location,
	)

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

func Indent(s, prefix string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return s
	}
	return regexp.MustCompile(`(?m)^`).ReplaceAllLiteralString(s, prefix)
}

func Dedent(s string) string {
	lines := strings.Split(s, "\n")
	minIndent := -1

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return s
	}

	var buffer bytes.Buffer
	for _, line := range lines {
		fmt.Fprintln(&buffer, strings.TrimPrefix(line, strings.Repeat(" ", minIndent)))
	}
	return strings.TrimSuffix(buffer.String(), "\n")
}

func Duration(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60

	var parts []string
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	if mins > 0 {
		if mins == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", mins))
		}
	}

	if len(parts) == 0 {
		return "0 minutes"
	}
	return strings.Join(parts, " ")
}
