package api

import (
	"context"
	"github.com/timwehrle/asana/internal/config"
	"net/url"
	"time"
)

type Workspace struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

func (w Workspace) ToYaml() config.DefaultWorkspace {
	return config.DefaultWorkspace{
		GID:  w.GID,
		Name: w.Name,
	}
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	endpoint := &url.URL{
		Path: "workspaces",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []Workspace `json:"data"`
	}

	if err := c.Response(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
