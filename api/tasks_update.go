package api

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

func (c *Client) MarkTaskAsDone(taskGID string) error {
	endpoint := &url.URL{
		Path: fmt.Sprintf("tasks/%s", taskGID),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"completed": true,
		},
	}

	resp, err := c.makeRequest(ctx, "PUT", endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to mark task %s as done: %w", taskGID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to mark task %s as done: status code %d", taskGID, resp.StatusCode)
	}
	return nil
}
