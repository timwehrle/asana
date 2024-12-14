package api

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

func (c *Client) GetTask(taskGID string) (*Task, error) {
	endpoint := &url.URL{
		Path: "tasks/" + taskGID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data Task `json:"data"`
	}

	if err := handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to handle response: %w", err)
	}

	return &result.Data, nil
}
