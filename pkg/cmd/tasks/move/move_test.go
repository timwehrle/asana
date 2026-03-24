package move

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

func TestRunMove_TaskToSectionByName(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	var requests []string
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.Path)

		switch req.Method + " " + req.URL.Path {
		case "GET /api/1.0/projects/P1/sections":
			return asana.MockResponse(http.StatusOK, []*asana.Section{
				{ID: "S1", SectionBase: asana.SectionBase{Name: "Backlog"}},
				{ID: "S2", SectionBase: asana.SectionBase{Name: "Ready"}},
			})
		case "POST /api/1.0/tasks/T1/addProject":
			body, err := (&asana.AssertRequest{Request: req}).Body()
			if err != nil {
				t.Fatalf("Body(): %v", err)
			}
			data := body["data"].(map[string]any)
			if got, want := data["project"], "P1"; got != want {
				t.Fatalf("project = %q; want %q", got, want)
			}
			if got, want := data["section"], "S1"; got != want {
				t.Fatalf("section = %q; want %q", got, want)
			}
			return asana.MockResponse(http.StatusOK, map[string]any{})
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return nil, nil
		}
	})

	opts := &MoveOptions{
		IO:          io,
		Prompter:    prompter.NewMockPrompter(),
		Client:      func() (*asana.Client, error) { return client, nil },
		TaskID:      "T1",
		ProjectID:   "P1",
		SectionName: "Backlog",
	}

	if err := runMove(opts); err != nil {
		t.Fatal(err)
	}

	wantRequests := []string{
		"GET /api/1.0/projects/P1/sections",
		"POST /api/1.0/tasks/T1/addProject",
	}
	if strings.Join(requests, "|") != strings.Join(wantRequests, "|") {
		t.Fatalf("requests = %v; want %v", requests, wantRequests)
	}
	if !strings.Contains(out.String(), "Task moved") {
		t.Fatalf("output = %q; want move confirmation", out.String())
	}
}
