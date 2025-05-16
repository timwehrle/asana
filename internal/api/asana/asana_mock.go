package asana

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MockClient struct {
	DoFunc   func(req *http.Request) (*http.Response, error)
	Requests []*http.Request
}

// NewMockClient creates a new mock client with the given response
func NewMockClient(status int, body any) (*MockClient, error) {
	responseBody, err := prepareBody(body)
	if err != nil {
		return nil, err
	}

	return &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}, nil
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, req)
	return m.DoFunc(req)
}

// MockResponse creates a mock HTTP response with the given status code and body
func MockResponse(status int, body any) (*http.Response, error) {
	var bodyContent []byte
	var err error

	switch v := body.(type) {
	case string:
		bodyContent = []byte(v)
	case []byte:
		bodyContent = v
	default:
		bodyContent, err = json.Marshal(map[string]any{
			"data": body,
		})
		if err != nil {
			return nil, err
		}
	}

	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBuffer(bodyContent)),
		Header:     make(http.Header),
	}, nil
}

// prepareBody prepares the body for repeated use
func prepareBody(body any) ([]byte, error) {
	switch v := body.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return json.Marshal(map[string]any{
			"data": body,
		})
	}
}

// AssertRequest provides helper methods to assert request properties
type AssertRequest struct {
	Request *http.Request
}

func (a *AssertRequest) Method() string {
	return a.Request.Method
}

func (a *AssertRequest) Path() string {
	return a.Request.URL.Path
}

func (a *AssertRequest) Query() url.Values {
	return a.Request.URL.Query()
}

func (a *AssertRequest) Header(key string) string {
	return a.Request.Header.Get(key)
}

func (a *AssertRequest) Body() (map[string]any, error) {
	if a.Request.Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(a.Request.Body)
	if err != nil {
		return nil, err
	}

	// Restore the body for subsequent reads
	a.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (a *AssertRequest) HasFeature(feature string) bool {
	enableHeader := a.Header("Asana-Enable")
	if enableHeader == "" {
		return false
	}
	features := strings.Split(enableHeader, ",")
	for _, f := range features {
		if f == feature {
			return true
		}
	}
	return false
}

// GetLastRequest returns an AssertRequest for the last request made to the mock client
func (m *MockClient) GetLastRequest() *AssertRequest {
	if len(m.Requests) == 0 {
		return nil
	}
	return &AssertRequest{Request: m.Requests[len(m.Requests)-1]}
}
