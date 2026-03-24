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
		for _, field := range []string{
			"name",
			"notes",
			"completed",
			"permalink_url",
			"dependencies.name",
			"dependencies.completed",
			"custom_fields.name",
			"custom_fields.resource_subtype",
			"custom_fields.text_value",
			"custom_fields.enum_value.name",
			"memberships.project.name",
			"memberships.section.name",
		} {
			if !strings.Contains(fields, field) {
				t.Fatalf("opt_fields = %q; missing %q", fields, field)
			}
		}

		return asana.MockResponse(http.StatusOK, map[string]any{
			"gid":           "T1",
			"name":          "Agent task",
			"notes":         "Branch: codex/foo\n\nImplement parser",
			"completed":     false,
			"due_on":        "2026-03-27",
			"permalink_url": "https://example.com/tasks/T1",
			"dependencies": []map[string]any{
				{
					"gid":       "D1",
					"name":      "Predecessor task",
					"completed": false,
				},
			},
			"custom_fields": []map[string]any{
				{
					"gid":              "CF1",
					"name":             "Priority",
					"resource_subtype": "enum",
					"enum_value": map[string]any{
						"gid":  "EV1",
						"name": "High",
					},
				},
				{
					"gid":              "CF2",
					"name":             "PR URL",
					"resource_subtype": "text",
					"text_value":       nil,
				},
			},
			"memberships": []map[string]any{
				{
					"project": map[string]any{
						"gid":  "P1",
						"name": "Website Tasks",
					},
					"section": map[string]any{
						"gid":  "S1",
						"name": "In Progress",
					},
				},
				{
					"project": map[string]any{
						"gid":  "P2",
						"name": "Agent Orchestration",
					},
					"section": map[string]any{
						"gid":  "S2",
						"name": "Ready for development",
					},
				},
			},
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
			DueOn        string `json:"due_on"`
			PermalinkURL string `json:"permalink_url"`
			Dependencies []struct {
				GID       string `json:"gid"`
				Name      string `json:"name"`
				Completed bool   `json:"completed"`
			} `json:"dependencies"`
			CustomFields []struct {
				GID       string  `json:"gid"`
				Name      string  `json:"name"`
				Type      string  `json:"type"`
				TextValue *string `json:"text_value"`
				EnumValue *struct {
					GID  string `json:"gid"`
					Name string `json:"name"`
				} `json:"enum_value"`
			} `json:"custom_fields"`
			Memberships []struct {
				Project struct {
					GID  string `json:"gid"`
					Name string `json:"name"`
				} `json:"project"`
				Section struct {
					GID  string `json:"gid"`
					Name string `json:"name"`
				} `json:"section"`
			} `json:"memberships"`
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
	if payload.Task.DueOn != "2026-03-27" {
		t.Fatalf("due_on = %q; want %q", payload.Task.DueOn, "2026-03-27")
	}
	if len(payload.Task.Dependencies) != 1 {
		t.Fatalf("dependencies length = %d; want 1", len(payload.Task.Dependencies))
	}
	if got, want := payload.Task.Dependencies[0].GID, "D1"; got != want {
		t.Fatalf("dependencies[0].gid = %q; want %q", got, want)
	}
	if payload.Task.Dependencies[0].Completed {
		t.Fatalf("dependencies[0].completed = true; want false")
	}
	if len(payload.Task.CustomFields) != 2 {
		t.Fatalf("custom_fields length = %d; want 2", len(payload.Task.CustomFields))
	}
	if got, want := payload.Task.CustomFields[0].Type, "enum"; got != want {
		t.Fatalf("custom_fields[0].type = %q; want %q", got, want)
	}
	if payload.Task.CustomFields[0].EnumValue == nil || payload.Task.CustomFields[0].EnumValue.Name != "High" {
		t.Fatalf("custom_fields[0].enum_value = %#v; want name High", payload.Task.CustomFields[0].EnumValue)
	}
	if payload.Task.CustomFields[1].TextValue != nil {
		t.Fatalf("custom_fields[1].text_value = %v; want nil", *payload.Task.CustomFields[1].TextValue)
	}
	if len(payload.Task.Memberships) != 2 {
		t.Fatalf("memberships length = %d; want 2", len(payload.Task.Memberships))
	}
	if got, want := payload.Task.Memberships[0].Project.GID, "P1"; got != want {
		t.Fatalf("memberships[0].project.gid = %q; want %q", got, want)
	}
	if got, want := payload.Task.Memberships[1].Section.Name, "Ready for development"; got != want {
		t.Fatalf("memberships[1].section.name = %q; want %q", got, want)
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
