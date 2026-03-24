package add

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

func TestRunAddComment_MentionCreator(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	oldDelay := followerPropagationDelay
	followerPropagationDelay = 0
	defer func() { followerPropagationDelay = oldDelay }()

	now := time.Date(2026, 3, 24, 9, 45, 0, 0, time.UTC)
	var requests []string
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.Path)

		switch req.Method + " " + req.URL.Path {
		case "GET /api/1.0/tasks/T1":
			fields := req.URL.Query().Get("opt_fields")
			for _, field := range []string{"created_by.name", "created_by.gid"} {
				if !strings.Contains(fields, field) {
					t.Fatalf("opt_fields = %q; missing %q", fields, field)
				}
			}
			return asana.MockResponse(http.StatusOK, map[string]any{
				"gid":        "T1",
				"name":       "Triage task",
				"created_by": map[string]any{"gid": "U1", "name": "Alice"},
			})
		case "POST /api/1.0/tasks/T1/addFollowers":
			body, err := (&asana.AssertRequest{Request: req}).Body()
			if err != nil {
				t.Fatalf("Body(): %v", err)
			}
			data := body["data"].(map[string]any)
			followers := data["followers"].([]any)
			if len(followers) != 1 || followers[0] != "U1" {
				t.Fatalf("followers = %v; want [U1]", followers)
			}
			return asana.MockResponse(http.StatusOK, map[string]any{})
		case "POST /api/1.0/tasks/T1/stories":
			body, err := (&asana.AssertRequest{Request: req}).Body()
			if err != nil {
				t.Fatalf("Body(): %v", err)
			}
			data := body["data"].(map[string]any)
			htmlText, _ := data["html_text"].(string)
			if !strings.Contains(htmlText, `data-asana-gid="U1"`) {
				t.Fatalf("html_text = %q; missing user mention", htmlText)
			}
			if !strings.Contains(htmlText, "Can you clarify the expected behavior?") {
				t.Fatalf("html_text = %q; missing comment text", htmlText)
			}

			return asana.MockResponse(http.StatusCreated, map[string]any{
				"gid":              "S1",
				"resource_subtype": "comment_added",
				"text":             "@Alice Can you clarify the expected behavior?",
				"html_text":        htmlText,
				"created_at":       now.Format(time.RFC3339),
				"created_by": map[string]any{
					"gid":  "ME",
					"name": "Me",
				},
			})
		default:
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
			return nil, nil
		}
	})

	opts := &AddOptions{
		IO:             io,
		Prompter:       prompter.NewMockPrompter(),
		Client:         func() (*asana.Client, error) { return client, nil },
		TaskID:         "T1",
		Text:           "Can you clarify the expected behavior?",
		MentionCreator: true,
		Output:         "json",
	}

	if err := runAdd(opts); err != nil {
		t.Fatal(err)
	}

	wantRequests := []string{
		"GET /api/1.0/tasks/T1",
		"POST /api/1.0/tasks/T1/addFollowers",
		"POST /api/1.0/tasks/T1/stories",
	}
	if strings.Join(requests, "|") != strings.Join(wantRequests, "|") {
		t.Fatalf("requests = %v; want %v", requests, wantRequests)
	}

	var payload struct {
		Story struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"story"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, out.String())
	}
	if payload.Story.ID != "S1" {
		t.Fatalf("story id = %q; want %q", payload.Story.ID, "S1")
	}
}
