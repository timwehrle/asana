package api

import (
	"net/url"

	"github.com/timwehrle/alfie/internal/workspace"
)

type Task struct {
	GID          string        `json:"gid"`
	Name         string        `json:"name"`
	DueOn        string        `json:"due_on"`
	CreatedBy    User          `json:"created_by"`
	HTMLNotes    string        `json:"html_notes"`
	Notes        string        `json:"notes"`
	Assignee     User          `json:"assignee"`
	Tags         []Tag         `json:"tags"`
	Link         string        `json:"permalink_url"`
	CustomFields []CustomField `json:"custom_fields"`
	Projects     []Project     `json:"projects"`
}

func (c *Client) GetTasks() ([]Task, error) {
	workspaceGID, _, err := workspace.LoadDefaultWorkspace()
	if err != nil {
		return nil, err
	}

	endpoint := &url.URL{
		Path: "tasks",
	}
	query := url.Values{}
	query.Set("workspace", workspaceGID)
	query.Set("opt_fields", "due_on,name,completed")
	query.Set("completed_since", "now")
	query.Set("assignee", "me")
	endpoint.RawQuery = query.Encode()

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []Task `json:"data"`
	}

	if err := handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
