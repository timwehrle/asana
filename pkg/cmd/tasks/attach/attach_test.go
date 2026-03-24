package attach

import (
	"net/http"
	"strings"
	"testing"

	"github.com/timwehrle/asana/internal/api/asana"
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

func TestNewCmdAttach_RunE_ParsesFlags(t *testing.T) {
	f, _, _ := factory.NewTestFactory()

	var sawOpts *AttachOptions
	cmd := NewCmdAttach(f, func(opts *AttachOptions) error {
		sawOpts = opts
		return nil
	})

	cmd.SetArgs([]string{
		"--task", "T1",
		"--url", "https://github.com/org/repo/pull/123",
		"--name", "PR #123: Short title",
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
	if sawOpts.URL != "https://github.com/org/repo/pull/123" {
		t.Fatalf("URL = %q; want PR URL", sawOpts.URL)
	}
	if sawOpts.Name != "PR #123: Short title" {
		t.Fatalf("Name = %q; want PR title", sawOpts.Name)
	}
}

func TestRunAttach_ExternalAttachment(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if got, want := req.Method, http.MethodPost; got != want {
			t.Fatalf("Method = %q; want %q", got, want)
		}
		if got, want := req.URL.Path, "/api/1.0/tasks/T1/attachments"; got != want {
			t.Fatalf("Path = %q; want %q", got, want)
		}

		body, err := (&asana.AssertRequest{Request: req}).Body()
		if err != nil {
			t.Fatalf("Body(): %v", err)
		}
		data := body["data"].(map[string]any)
		if got, want := data["resource_subtype"], "external"; got != want {
			t.Fatalf("resource_subtype = %q; want %q", got, want)
		}
		if got, want := data["url"], "https://github.com/org/repo/pull/123"; got != want {
			t.Fatalf("url = %q; want %q", got, want)
		}
		if got, want := data["name"], "PR #123: Short title"; got != want {
			t.Fatalf("name = %q; want %q", got, want)
		}

		return asana.MockResponse(http.StatusCreated, map[string]any{
			"gid":              "A1",
			"name":             "PR #123: Short title",
			"resource_subtype": "external",
			"permanent_url":    "https://app.asana.com/0/0/A1/f",
			"view_url":         "https://github.com/org/repo/pull/123",
			"connected_to_app": false,
		})
	})

	opts := &AttachOptions{
		IO:     io,
		Client: func() (*asana.Client, error) { return client, nil },
		TaskID: "T1",
		URL:    "https://github.com/org/repo/pull/123",
		Name:   "PR #123: Short title",
	}

	if err := runAttach(opts); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Attached PR #123: Short title to task T1") {
		t.Fatalf("output = %q; want attachment confirmation", out.String())
	}
}
