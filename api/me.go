package api

type User struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

func (c *Client) Me() (User, error) {
	resp, err := c.makeRequest("GET", "/users/me", nil)
	if err != nil {
		return User{}, err
	}

	var result struct {
		Data User `json:"data"`
	}

	if err := handleResponse(resp, &result); err != nil {
		return User{}, err
	}

	return result.Data, nil
}
