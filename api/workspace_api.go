package api

import (
	"context"
	"net/url"
	"time"
)

type Workspace struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	endpoint := &url.URL{
		Path: "workspaces",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []Workspace `json:"data"`
	}

	if err := handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
