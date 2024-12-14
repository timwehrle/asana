package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// New creates a new API client with the provided token.
func New(token string) *Client {
	return &Client{
		BaseURL:    "https://api.asana.com/api/1.0/",
		Token:      token,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// makeRequest sends an HTTP request to the Asana API and returns the response.
func (c *Client) makeRequest(ctx context.Context, method string, endpoint *url.URL, body any) (*http.Response, error) {
	fullURL, err := c.buildFullURL(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build full URL: %w", err)
	}

	reqBody, err := c.marshalBody(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func handleResponse(resp *http.Response, result any) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("API request failed with status: %s, but failed to read response body: %w", resp.Status, err)
		}
		return fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) buildFullURL(endpoint *url.URL) (string, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
	}
	return baseURL.ResolveReference(endpoint).String(), nil
}

func (c *Client) marshalBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	return json.Marshal(body)
}
