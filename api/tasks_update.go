package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) UpdateTask(taskGID string, updates map[string]any) error {
	endpoint := &url.URL{
		Path: fmt.Sprintf("tasks/%s", taskGID),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload := map[string]any{
		"data": updates,
	}

	resp, err := c.Request(ctx, "PUT", endpoint, payload)
	if err != nil {
		return fmt.Errorf("failed to mark task %s: %w", taskGID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mark task %s: status code %d", taskGID, resp.StatusCode)
	}

	return nil
}
