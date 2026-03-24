package view

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
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

func TestNewCmdView_RunE_ParsesTaskAndOutputFlags(t *testing.T) {
	f, _, _ := factory.NewTestFactory()

	var sawOpts *ViewOptions
	cmd := NewCmdView(f, func(opts *ViewOptions) error {
		sawOpts = opts
		return nil
	})

	cmd.SetArgs([]string{"--task", "T1", "--output", "json"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if sawOpts == nil {
		t.Fatal("runF was never called")
	}
	if sawOpts.TaskID != "T1" {
		t.Fatalf("TaskID = %q; want %q", sawOpts.TaskID, "T1")
	}
	if sawOpts.Output != "json" {
		t.Fatalf("Output = %q; want %q", sawOpts.Output, "json")
	}
}

func TestViewRun_JSONByTaskID(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if got, want := req.Method, http.MethodGet; got != want {
			t.Fatalf("Method = %q; want %q", got, want)
		}
		if got, want := req.URL.Path, "/api/1.0/tasks/T1"; got != want {
			t.Fatalf("Path = %q; want %q", got, want)
		}
		fields := req.URL.Query().Get("opt_fields")
		for _, field := range []string{"name", "notes", "completed", "permalink_url"} {
			if !strings.Contains(fields, field) {
				t.Fatalf("opt_fields = %q; missing %q", fields, field)
			}
		}

		return asana.MockResponse(http.StatusOK, map[string]any{
			"gid":           "T1",
			"name":          "Agent task",
			"notes":         "Branch: codex/foo\n\nImplement parser",
			"completed":     false,
			"permalink_url": "https://example.com/tasks/T1",
		})
	})

	opts := &ViewOptions{
		IO:       io,
		Prompter: prompter.NewMockPrompter(),
		Config: func() (*config.Config, error) {
			t.Fatal("Config should not be called when --task is provided")
			return nil, nil
		},
		Client: func() (*asana.Client, error) { return client, nil },
		TaskID: "T1",
		Output: "json",
	}

	if err := viewRun(opts); err != nil {
		t.Fatal(err)
	}

	var payload struct {
		Task struct {
			GID          string `json:"gid"`
			Notes        string `json:"notes"`
			PermalinkURL string `json:"permalink_url"`
		} `json:"task"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, out.String())
	}
	if payload.Task.GID != "T1" {
		t.Fatalf("gid = %q; want %q", payload.Task.GID, "T1")
	}
	if !strings.Contains(payload.Task.Notes, "Branch: codex/foo") {
		t.Fatalf("notes = %q; want branch header", payload.Task.Notes)
	}
}

func TestViewRun_NonTTYRequiresTaskID(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.IsStdinTTY = false

	opts := &ViewOptions{
		IO:       io,
		Prompter: prompter.NewMockPrompter(),
		Config: func() (*config.Config, error) {
			t.Fatal("Config should not be called when prompting is disallowed")
			return nil, nil
		},
		Client: func() (*asana.Client, error) {
			t.Fatal("Client should not be called when prompting is disallowed")
			return nil, nil
		},
	}

	err := viewRun(opts)
	if err == nil || !strings.Contains(err.Error(), "--task") {
		t.Fatalf("expected --task guidance error, got %v", err)
	}
}
