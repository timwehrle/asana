package api

import "fmt"

func (c *Client) GetTask(taskGID string) (*Task, error) {
	resp, err := c.makeRequest("GET", "/tasks/"+taskGID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data Task `json:"data"`
	}

	if err := handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to handle response: %v", err)
	}

	return &result.Data, nil
}
