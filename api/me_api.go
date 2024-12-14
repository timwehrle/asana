package api

import (
	"context"
	"net/url"
	"time"
)

func (c *Client) GetMe() (User, error) {
	endpoint := &url.URL{
		Path: "users/me",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Data User `json:"data"`
	}

	if err := handleResponse(resp, &result); err != nil {
		return User{}, err
	}

	return result.Data, nil
}
