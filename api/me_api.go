package api

import (
	"context"
	"net/url"
	"time"
)

type User struct {
	ID   string `json:"gid"`
	Name string `json:"name"`
}

func (c *Client) GetMe() (User, error) {
	endpoint := &url.URL{
		Path: "users/me",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Data User `json:"data"`
	}

	if err := c.Response(resp, &result); err != nil {
		return User{}, err
	}

	return result.Data, nil
}

func (u User) Username() string {
	return u.Name
}

func (u User) GID() string {
	return u.ID
}
