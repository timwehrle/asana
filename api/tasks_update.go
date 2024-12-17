package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) MarkTaskAsDone(taskGID string) error {
	endpoint := &url.URL{
		Path: fmt.Sprintf("tasks/%s", taskGID),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload := map[string]any{
		"data": map[string]any{
			"completed": true,
		},
	}

	resp, err := c.Request(ctx, "PUT", endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to mark task %s as done: %w", taskGID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mark task %s as done: status code %d", taskGID, resp.StatusCode)
	}
	return nil
}
