package update

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
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

func TestNewCmdUpdate_RunE_ParsesAutomationFlags(t *testing.T) {
	f, _, _ := factory.NewTestFactory()

	var sawOpts *UpdateOptions
	cmd := NewCmdUpdate(f, func(opts *UpdateOptions) error {
		sawOpts = opts
		return nil
	})

	cmd.SetArgs([]string{
		"--task", "T1",
		"--prepend-notes", "Branch: codex/foo\n\n",
		"--output", "json",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if sawOpts == nil {
		t.Fatal("runF was never called")
	}
	if sawOpts.TaskID != "T1" {
		t.Fatalf("TaskID = %q; want %q", sawOpts.TaskID, "T1")
	}
	if sawOpts.PrependNotes != "Branch: codex/foo\n\n" {
		t.Fatalf("PrependNotes = %q; want branch header", sawOpts.PrependNotes)
	}
	if sawOpts.Output != "json" {
		t.Fatalf("Output = %q; want %q", sawOpts.Output, "json")
	}
}

func TestRunUpdate_PrependNotes_JSONOutput(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	var putSeen bool
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/api/1.0/tasks/T1":
			return asana.MockResponse(http.StatusOK, map[string]any{
				"gid":           "T1",
				"name":          "Agent task",
				"notes":         "Implement parser",
				"completed":     false,
				"permalink_url": "https://example.com/tasks/T1",
			})
		case req.Method == http.MethodPut && req.URL.Path == "/api/1.0/tasks/T1":
			putSeen = true

			body := &asana.AssertRequest{Request: req}
			payload, err := body.Body()
			if err != nil {
				t.Fatalf("Body(): %v", err)
			}
			data := payload["data"].(map[string]any)
			if got, want := data["notes"], "Branch: codex/foo\n\nImplement parser"; got != want {
				t.Fatalf("notes payload = %q; want %q", got, want)
			}

			return asana.MockResponse(http.StatusOK, map[string]any{
				"gid":           "T1",
				"name":          "Agent task",
				"notes":         "Branch: codex/foo\n\nImplement parser",
				"completed":     false,
				"permalink_url": "https://example.com/tasks/T1",
			})
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
			return nil, nil
		}
	})

	opts := &UpdateOptions{
		IO:       io,
		Prompter: prompter.NewMockPrompter(),
		Config: func() (*config.Config, error) {
			t.Fatal("Config should not be called when --task is provided")
			return nil, nil
		},
		Client:       func() (*asana.Client, error) { return client, nil },
		TaskID:       "T1",
		PrependNotes: "Branch: codex/foo\n\n",
		Output:       "json",
	}

	if err := runUpdate(opts); err != nil {
		t.Fatal(err)
	}
	if !putSeen {
		t.Fatal("expected PUT request")
	}

	var payload struct {
		Task struct {
			GID   string `json:"gid"`
			Notes string `json:"notes"`
		} `json:"task"`
		UpdatedFields []string `json:"updated_fields"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, out.String())
	}
	if payload.Task.GID != "T1" {
		t.Fatalf("gid = %q; want %q", payload.Task.GID, "T1")
	}
	if payload.Task.Notes != "Branch: codex/foo\n\nImplement parser" {
		t.Fatalf("notes = %q; want prepended notes", payload.Task.Notes)
	}
	if len(payload.UpdatedFields) == 0 || payload.UpdatedFields[0] != "notes" {
		t.Fatalf("updated_fields = %v; want notes", payload.UpdatedFields)
	}
}

func TestRunUpdate_NotesFileReplacesNotes(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	dir := t.TempDir()
	notesPath := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(notesPath, []byte("File based notes"), 0o600); err != nil {
		t.Fatal(err)
	}

	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/api/1.0/tasks/T1":
			return asana.MockResponse(http.StatusOK, map[string]any{
				"gid":   "T1",
				"name":  "Agent task",
				"notes": "Existing notes",
			})
		case req.Method == http.MethodPut && req.URL.Path == "/api/1.0/tasks/T1":
			payload, err := (&asana.AssertRequest{Request: req}).Body()
			if err != nil {
				t.Fatalf("Body(): %v", err)
			}
			data := payload["data"].(map[string]any)
			if got, want := data["notes"], "File based notes"; got != want {
				t.Fatalf("notes payload = %q; want %q", got, want)
			}
			return asana.MockResponse(http.StatusOK, map[string]any{
				"gid":   "T1",
				"name":  "Agent task",
				"notes": "File based notes",
			})
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
			return nil, nil
		}
	})

	opts := &UpdateOptions{
		IO:        io,
		Prompter:  prompter.NewMockPrompter(),
		Client:    func() (*asana.Client, error) { return client, nil },
		TaskID:    "T1",
		NotesFile: notesPath,
	}

	if err := runUpdate(opts); err != nil {
		t.Fatal(err)
	}
}

func TestRunUpdate_NonTTYRequiresTaskAndAction(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	io.IsStdinTTY = false

	opts := &UpdateOptions{
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

	err := runUpdate(opts)
	if err == nil || !strings.Contains(err.Error(), "--task") {
		t.Fatalf("expected --task guidance error, got %v", err)
	}
}

func TestRunUpdate_RejectsConflictingNotesFlags(t *testing.T) {
	io, _, _, _ := iostreams.Test()

	opts := &UpdateOptions{
		IO:           io,
		Prompter:     prompter.NewMockPrompter(),
		TaskID:       "T1",
		Notes:        "replace",
		PrependNotes: "prepend",
		Client: func() (*asana.Client, error) {
			t.Fatal("Client should not be called for invalid flag combinations")
			return nil, nil
		},
	}

	err := runUpdate(opts)
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("expected mutually exclusive error, got %v", err)
	}
}
