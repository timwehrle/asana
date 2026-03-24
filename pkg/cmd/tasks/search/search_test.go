package search

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

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

func TestNewCmdSearch_RunE_ParsesAutomationFlags(t *testing.T) {
	f, _, _ := factory.NewTestFactory()

	var sawOpts *SearchOptions
	cmd := NewCmdSearch(f, func(opts *SearchOptions) error {
		sawOpts = opts
		return nil
	})

	cmd.SetArgs([]string{
		"--query", "Agents",
		"--output", "json",
		"--incomplete",
		"--limit", "5",
		"--page-offset", "CURSOR-1",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if sawOpts == nil {
		t.Fatal("runF was never called")
	}
	if sawOpts.Output != "json" {
		t.Fatalf("Output = %q; want %q", sawOpts.Output, "json")
	}
	if !sawOpts.IncompleteOnly {
		t.Fatal("IncompleteOnly = false; want true")
	}
	if sawOpts.Limit != 5 {
		t.Fatalf("Limit = %d; want 5", sawOpts.Limit)
	}
	if sawOpts.PageOffset != "CURSOR-1" {
		t.Fatalf("PageOffset = %q; want %q", sawOpts.PageOffset, "CURSOR-1")
	}
}

func TestRunSearch_JSONOutputUsesAutomationFilters(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	var requestCount int
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		requestCount++

		if got, want := req.Method, http.MethodGet; got != want {
			t.Fatalf("Method = %q; want %q", got, want)
		}
		if got, want := req.URL.Path, "/api/1.0/workspaces/W123/tasks/search"; got != want {
			t.Fatalf("Path = %q; want %q", got, want)
		}

		query := req.URL.Query()
		if got, want := query.Get("text"), "Agents"; got != want {
			t.Fatalf("text = %q; want %q", got, want)
		}
		if got, want := query.Get("completed"), "false"; got != want {
			t.Fatalf("completed = %q; want %q", got, want)
		}
		if got, want := query.Get("limit"), "5"; got != want {
			t.Fatalf("limit = %q; want %q", got, want)
		}
		if got, want := query.Get("offset"), "CURSOR-1"; got != want {
			t.Fatalf("offset = %q; want %q", got, want)
		}

		fields := query.Get("opt_fields")
		for _, field := range []string{"name", "completed", "due_on", "permalink_url"} {
			if !strings.Contains(fields, field) {
				t.Fatalf("opt_fields = %q; missing %q", fields, field)
			}
		}

		return asana.MockResponse(http.StatusOK, `{
			"data":[
				{
					"gid":"T1",
					"name":"Agent queue task",
					"completed":false,
					"due_on":"2026-03-20",
					"permalink_url":"https://example.com/tasks/T1"
				}
			],
			"next_page":{
				"offset":"CURSOR-2",
				"path":"/workspaces/W123/tasks/search?offset=CURSOR-2",
				"uri":"https://example.com/next"
			}
		}`)
	})

	opts := &SearchOptions{
		IO: io,
		Config: func() (*config.Config, error) {
			return &config.Config{Workspace: &asana.Workspace{ID: "W123"}}, nil
		},
		Client:         func() (*asana.Client, error) { return client, nil },
		Query:          "Agents",
		Output:         "json",
		IncompleteOnly: true,
		Limit:          5,
		PageOffset:     "CURSOR-1",
	}

	if err := runSearch(opts); err != nil {
		t.Fatal(err)
	}

	if requestCount != 1 {
		t.Fatalf("requestCount = %d; want 1", requestCount)
	}

	var payload struct {
		Tasks []struct {
			GID          string `json:"gid"`
			Name         string `json:"name"`
			Completed    bool   `json:"completed"`
			PermalinkURL string `json:"permalink_url"`
		} `json:"tasks"`
		NextPageOffset string `json:"next_page_offset"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, out.String())
	}

	if len(payload.Tasks) != 1 {
		t.Fatalf("len(tasks) = %d; want 1", len(payload.Tasks))
	}
	if payload.Tasks[0].GID != "T1" {
		t.Fatalf("gid = %q; want %q", payload.Tasks[0].GID, "T1")
	}
	if payload.NextPageOffset != "CURSOR-2" {
		t.Fatalf("next_page_offset = %q; want %q", payload.NextPageOffset, "CURSOR-2")
	}
}
