package list

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

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

func TestRunListComments_JSONByTaskID(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	now := time.Date(2026, 3, 24, 9, 30, 0, 0, time.UTC)
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if got, want := req.Method, http.MethodGet; got != want {
			t.Fatalf("Method = %q; want %q", got, want)
		}
		if got, want := req.URL.Path, "/api/1.0/tasks/T1/stories"; got != want {
			t.Fatalf("Path = %q; want %q", got, want)
		}

		fields := req.URL.Query().Get("opt_fields")
		for _, field := range []string{"created_by.name", "created_by.gid", "created_at", "text", "html_text", "resource_subtype"} {
			if !strings.Contains(fields, field) {
				t.Fatalf("opt_fields = %q; missing %q", fields, field)
			}
		}

		return asana.MockResponse(http.StatusOK, []*asana.Story{
			{
				ID:              "S1",
				ResourceSubtype: "comment_added",
				StoryBase: asana.StoryBase{
					Text:     "Need a more detailed spec here.",
					HTMLText: "<body>Need a more detailed spec here.</body>",
				},
				CreatedAt: &now,
				CreatedBy: &asana.User{ID: "U1", Name: "Alice"},
			},
			{
				ID:              "S2",
				ResourceSubtype: "assigned",
				StoryBase: asana.StoryBase{
					Text: "Assigned to Bob",
				},
				CreatedAt: &now,
				CreatedBy: &asana.User{ID: "U2", Name: "System"},
			},
		})
	})

	opts := &ListOptions{
		IO:       io,
		Prompter: prompter.NewMockPrompter(),
		Client:   func() (*asana.Client, error) { return client, nil },
		TaskID:   "T1",
		Output:   "json",
	}

	if err := runList(opts); err != nil {
		t.Fatal(err)
	}

	var payload struct {
		TaskID   string `json:"task_id"`
		Comments []struct {
			ID        string `json:"id"`
			Text      string `json:"text"`
			CreatedBy struct {
				GID  string `json:"gid"`
				Name string `json:"name"`
			} `json:"created_by"`
		} `json:"comments"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, out.String())
	}

	if payload.TaskID != "T1" {
		t.Fatalf("TaskID = %q; want %q", payload.TaskID, "T1")
	}
	if len(payload.Comments) != 1 {
		t.Fatalf("len(comments) = %d; want 1", len(payload.Comments))
	}
	if payload.Comments[0].ID != "S1" {
		t.Fatalf("comment id = %q; want %q", payload.Comments[0].ID, "S1")
	}
}
