package create

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (fn transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func newTestClient(fn transportFunc) *asana.Client {
	return asana.NewClient(&http.Client{Transport: fn})
}

func TestNewCmdCreate_RunE(t *testing.T) {
	f, _, _ := factory.NewTestFactory()

	var sawOpts *CreateOptions
	cmd := NewCmdCreate(f, func(opts *CreateOptions) error {
		sawOpts = opts
		return nil
	})

	cmd.SetArgs([]string{
		"--name", "My Task",
		"--assignee", "me",
		"--due", "2025-01-01",
		"--description", "Test description",
		"--parent", "T123",
		"--project", "P123",
		"--section-name", "Backlog",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if sawOpts == nil {
		t.Fatal("runF was never called")
	}

	if sawOpts.Name != "My Task" {
		t.Errorf("Name = %q; want %q", sawOpts.Name, "My Task")
	}
	if sawOpts.Assignee != "me" {
		t.Errorf("Assignee = %q; want %q", sawOpts.Assignee, "me")
	}
	if sawOpts.Due != "2025-01-01" {
		t.Errorf("Due = %q; want %q", sawOpts.Due, "2025-01-01")
	}
	if sawOpts.Description != "Test description" {
		t.Errorf("Description = %q; want %q", sawOpts.Description, "Test description")
	}
	if sawOpts.ParentID != "T123" {
		t.Errorf("ParentID = %q; want %q", sawOpts.ParentID, "T123")
	}
	if sawOpts.ProjectID != "P123" {
		t.Errorf("ProjectID = %q; want %q", sawOpts.ProjectID, "P123")
	}
	if sawOpts.SectionName != "Backlog" {
		t.Errorf("SectionName = %q; want %q", sawOpts.SectionName, "Backlog")
	}
}

func TestRunCreate_ConfigError(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	opts := &CreateOptions{
		IO: io,
		Config: func() (*config.Config, error) {
			return nil, errors.New("no config")
		},
		Client: func() (*asana.Client, error) {
			return nil, nil
		},
	}

	err := runCreate(opts)
	if err == nil || !strings.Contains(err.Error(), "failed to load config") {
		t.Fatalf("expected config error, got %v", err)
	}
}

func TestGetOrPromptDueDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		input   string
		wantDay string
	}{
		{
			name:    "today",
			input:   "today",
			wantDay: now.Format(time.DateOnly),
		},
		{
			name:    "tomorrow",
			input:   "tomorrow",
			wantDay: now.AddDate(0, 0, 1).Format(time.DateOnly),
		},
		{
			name:    "explicit date",
			input:   "2025-01-10",
			wantDay: "2025-01-10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &CreateOptions{Due: tt.input}

			got, err := getOrPromptDueDate(opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("got nil date")
			}

			gotDay := time.Time(*got).Format(time.DateOnly)
			if gotDay != tt.wantDay {
				t.Fatalf("date = %v; want %v", gotDay, tt.wantDay)
			}
		})
	}
}

func TestGetOrPromptDueDate_Invalid(t *testing.T) {
	opts := &CreateOptions{Due: "not-a-date"}

	_, err := getOrPromptDueDate(opts)
	if err == nil || !strings.Contains(err.Error(), "invalid due date") {
		t.Fatalf("expected invalid-date error, got %v", err)
	}
}

func TestRunCreate_ParentOnlySkipsProjectPrompts(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	var createBody map[string]any
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		switch req.Method + " " + req.URL.Path {
		case "GET /api/1.0/users":
			return asana.MockResponse(http.StatusOK, []*asana.User{
				{ID: "U1", Name: "Alice"},
			})
		case "POST /api/1.0/tasks":
			payload, err := (&asana.AssertRequest{Request: req}).Body()
			if err != nil {
				t.Fatalf("Body(): %v", err)
			}
			createBody = payload["data"].(map[string]any)
			return asana.MockResponse(http.StatusCreated, map[string]any{
				"gid":           "T999",
				"name":          "Subtask",
				"permalink_url": "https://example.com/tasks/T999",
			})
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return nil, nil
		}
	})

	opts := &CreateOptions{
		IO: io,
		Config: func() (*config.Config, error) {
			return &config.Config{Workspace: &asana.Workspace{ID: "W123"}, UserID: "U1"}, nil
		},
		Client:      func() (*asana.Client, error) { return client, nil },
		Name:        "Subtask",
		Assignee:    "me",
		Due:         "2026-03-24",
		Description: "Details",
		ParentID:    "PARENT1",
	}

	if err := runCreate(opts); err != nil {
		t.Fatal(err)
	}

	if got, want := createBody["parent"], "PARENT1"; got != want {
		t.Fatalf("parent = %q; want %q", got, want)
	}
	if _, ok := createBody["projects"]; ok {
		t.Fatalf("projects unexpectedly set: %v", createBody["projects"])
	}
	if _, ok := createBody["memberships"]; ok {
		t.Fatalf("memberships unexpectedly set: %v", createBody["memberships"])
	}
}

func TestRunCreate_ProjectAndSectionNameResolveMembership(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	var createBody map[string]any
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		switch req.Method + " " + req.URL.Path {
		case "GET /api/1.0/users":
			return asana.MockResponse(http.StatusOK, []*asana.User{
				{ID: "U1", Name: "Alice"},
			})
		case "GET /api/1.0/projects/P123/sections":
			return asana.MockResponse(http.StatusOK, []*asana.Section{
				{ID: "S1", SectionBase: asana.SectionBase{Name: "Backlog"}},
				{ID: "S2", SectionBase: asana.SectionBase{Name: "Ready"}},
			})
		case "POST /api/1.0/tasks":
			payload, err := (&asana.AssertRequest{Request: req}).Body()
			if err != nil {
				t.Fatalf("Body(): %v", err)
			}
			createBody = payload["data"].(map[string]any)
			return asana.MockResponse(http.StatusCreated, map[string]any{
				"gid":           "T999",
				"name":          "Scoped task",
				"permalink_url": "https://example.com/tasks/T999",
			})
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return nil, nil
		}
	})

	opts := &CreateOptions{
		IO: io,
		Config: func() (*config.Config, error) {
			return &config.Config{Workspace: &asana.Workspace{ID: "W123"}, UserID: "U1"}, nil
		},
		Client:      func() (*asana.Client, error) { return client, nil },
		Name:        "Scoped task",
		Assignee:    "me",
		Due:         "2026-03-24",
		Description: "Details",
		ProjectID:   "P123",
		SectionName: "Backlog",
	}

	if err := runCreate(opts); err != nil {
		t.Fatal(err)
	}

	projects := createBody["projects"].([]any)
	if len(projects) != 1 || projects[0] != "P123" {
		t.Fatalf("projects = %v; want [P123]", projects)
	}
	memberships := createBody["memberships"].([]any)
	if len(memberships) != 1 {
		t.Fatalf("memberships = %v; want length 1", memberships)
	}
	membership := memberships[0].(map[string]any)
	if got, want := membership["project"], "P123"; got != want {
		t.Fatalf("membership.project = %q; want %q", got, want)
	}
	if got, want := membership["section"], "S1"; got != want {
		t.Fatalf("membership.section = %q; want %q", got, want)
	}
}

func TestRunCreate_SectionNameRequiresProject(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	opts := &CreateOptions{
		IO:          io,
		Config:      func() (*config.Config, error) { return &config.Config{Workspace: &asana.Workspace{ID: "W123"}}, nil },
		Client:      func() (*asana.Client, error) { return &asana.Client{}, nil },
		Name:        "Task",
		Assignee:    "me",
		Due:         "2026-03-24",
		SectionName: "Backlog",
	}

	err := runCreate(opts)
	if err == nil || !strings.Contains(err.Error(), "--project") {
		t.Fatalf("expected --project guidance error, got %v", err)
	}
}
