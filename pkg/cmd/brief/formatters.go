package brief

import (
	"fmt"
	"github.com/timwehrle/alfie/api"
	"github.com/timwehrle/alfie/utils"
	"strings"
)

func formatNames(tasks []api.Task) []string {
	names := make([]string, len(tasks))
	for i, task := range tasks {
		names[i] = fmt.Sprintf("[%s] %s", utils.FormatDate(task.DueOn), task.Name)
	}
	return names
}

func formatList(prefix string, items []string) string {
	if len(items) == 0 {
		return prefix + "None"
	}
	return prefix + strings.Join(items, ", ")
}

func formatProjects(projects []api.Project) string {
	names := make([]string, len(projects))
	for i, project := range projects {
		names[i] = project.Name
	}
	return formatList("Projects: ", names)
}

func formatTags(tags []api.Tag) string {
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	return formatList("Tags: ", names)
}

func formatNotes(notes string) string {
	if notes == "" {
		return ""
	}
	return utils.BoldUnderline.Sprintf("Description:") + "\n" + notes + "\n"
}
