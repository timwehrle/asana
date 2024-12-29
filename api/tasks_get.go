package api

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/timwehrle/asana/internal/config"
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
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	endpoint := &url.URL{
		Path: "tasks",
	}
	query := url.Values{}
	query.Set("config", config.Workspace.GID)
	query.Set("opt_fields", "due_on,name,completed")
	query.Set("completed_since", "now")
	query.Set("assignee", "me")
	endpoint.RawQuery = query.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []Task `json:"data"`
	}

	if err := c.Response(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) GetTask(taskGID string) (*Task, error) {
	endpoint := &url.URL{
		Path: "tasks/" + taskGID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data Task `json:"data"`
	}

	if err := c.Response(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to handle response: %w", err)
	}

	return &result.Data, nil
}
