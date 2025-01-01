package api

import "fmt"

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("API error with status code %d: %s", e.StatusCode, e.Message)
}

type Tag struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type CustomField struct {
	GID       string `json:"gid"`
	CreatedBy User   `json:"created_by"`
}

type Membership struct {
	Project Project `json:"project"`
	Section Section `json:"section"`
}

type Section struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

type Project struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}
