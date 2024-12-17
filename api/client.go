package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	httpClient *http.Client
	once       sync.Once
}

// New creates a new API client with the provided token
func New(token string) *Client {
	return &Client{
		BaseURL: "https://api.asana.com/api/1.0/",
		Token:   token,
	}
}

// WithHTTPCLient allows customizing the HTTP client
func (c *Client) WithHTTPCLient(client *http.Client) *Client {
	c.httpClient = client
	return c
}

// getHTTPClient ensures a default client is created if not provided
func (c *Client) getHTTPClient() *http.Client {
	c.once.Do(func() {
		if c.httpClient == nil {
			c.httpClient = &http.Client{
				Timeout: 30 * time.Second,
				Transport: &http.Transport{
					MaxIdleConns:        10,
					MaxIdleConnsPerHost: 5,
					IdleConnTimeout:     90 * time.Second,
				},
			}
		}
	})
	return c.httpClient
}

// Request sends an HTTP request to the Asana API and returns the response
func (c *Client) Request(ctx context.Context, method string, endpoint *url.URL, body any) (*http.Response, error) {
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
	req.Header.Set("Accept", "application/json")

	return c.getHTTPClient().Do(req)
}

// Response handles the HTTP response, checking status and decoding result
func (c *Client) Response(resp *http.Response, result any) (err error) {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(resp)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// handleErrorResponse provides detailed error information
func (c *Client) handleErrorResponse(resp *http.Response) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("API request failed with status %s, and error reading response body: %w", resp.Status, err)
	}

	return fmt.Errorf("API request failed with status %s, body: %s", resp.Status, string(bodyBytes))
}

// buildFullURL constructs the full URL for the API request
func (c *Client) buildFullURL(endpoint *url.URL) (string, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
	}
	return baseURL.ResolveReference(endpoint).String(), nil
}

// marshalBody converts the request body to JSON
func (c *Client) marshalBody(body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	return json.Marshal(body)
}
