package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
}

func New(token string) *Client {
	return &Client{
		BaseURL: "https://api.asana.com/api/1.0",
		Token:   token,
	}
}

func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := c.BaseURL + endpoint

	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func handleResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"API request failed with status: %s, body: %s",
			resp.Status, string(bodyBytes),
		)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
