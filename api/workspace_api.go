package api

type Workspace struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	resp, err := c.makeRequest("GET", "/workspaces", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []Workspace `json:"data"`
	}

	if err := handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
