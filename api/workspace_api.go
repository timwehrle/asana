package api

import "net/url"

type Workspace struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	endpoint := &url.URL{
		Path: "workspaces",
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
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
