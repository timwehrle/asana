package api

import "net/url"

func (c *Client) GetMe() (User, error) {
	endpoint := &url.URL{
		Path: "users/me",
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
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
