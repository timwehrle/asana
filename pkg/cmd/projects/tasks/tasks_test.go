package tasks

import (
	"net/http"
	"strings"
	"testing"

	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (fn transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func newTestClient(fn transportFunc) *asana.Client {
	return asana.NewClient(&http.Client{Transport: fn})
}

func TestRunTasks_ByProjectIDWithSections(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	var requests []string
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.Path)

		switch req.Method + " " + req.URL.Path {
		case "GET /api/1.0/projects/P1":
			return asana.MockResponse(http.StatusOK, map[string]any{
				"gid":  "P1",
				"name": "Triage Project",
			})
		case "GET /api/1.0/projects/P1/sections":
			return asana.MockResponse(http.StatusOK, []*asana.Section{
				{ID: "S1", SectionBase: asana.SectionBase{Name: "Backlog"}},
				{ID: "S2", SectionBase: asana.SectionBase{Name: "Mother tasks"}},
			})
		case "GET /api/1.0/sections/S1/tasks":
			return asana.MockResponse(http.StatusOK, []*asana.Task{
				{ID: "T1", TaskBase: asana.TaskBase{Name: "First task"}},
				{ID: "T2", TaskBase: asana.TaskBase{Name: "Second task"}},
			})
		case "GET /api/1.0/sections/S2/tasks":
			return asana.MockResponse(http.StatusOK, []*asana.Task{
				{ID: "T3", TaskBase: asana.TaskBase{Name: "Third task"}},
			})
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return nil, nil
		}
	})

	opts := &TasksOptions{
		IO:           io,
		Prompter:     prompter.NewMockPrompter(),
		Client:       func() (*asana.Client, error) { return client, nil },
		ID:           "P1",
		WithSections: true,
	}

	if err := runTasks(opts); err != nil {
		t.Fatal(err)
	}

	wantRequests := []string{
		"GET /api/1.0/projects/P1",
		"GET /api/1.0/projects/P1/sections",
		"GET /api/1.0/sections/S1/tasks",
		"GET /api/1.0/sections/S2/tasks",
	}
	if strings.Join(requests, "|") != strings.Join(wantRequests, "|") {
		t.Fatalf("requests = %v; want %v", requests, wantRequests)
	}

	output := out.String()
	for _, want := range []string{
		"=== Backlog (id: S1) ===",
		"1. First task (id: T1)",
		"2. Second task (id: T2)",
		"=== Mother tasks (id: S2) ===",
		"3. Third task (id: T3)",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestRunTasks_NonTTYRequiresID(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.IsStdinTTY = false

	err := runTasks(&TasksOptions{
		IO:       io,
		Prompter: prompter.NewMockPrompter(),
		Client: func() (*asana.Client, error) {
			t.Fatal("Client should not be called when prompting is disallowed")
			return nil, nil
		},
	})
	if err == nil || !strings.Contains(err.Error(), "--id") {
		t.Fatalf("expected --id guidance error, got %v", err)
	}
}
