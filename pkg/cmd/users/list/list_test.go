package list

import (
	"errors"
	"net/http"
	"regexp"
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

func newTestClient(mock *asana.MockClient) *asana.Client {
	httpClient := &http.Client{
		Transport: transportFunc(mock.Do),
	}
	return asana.NewClient(httpClient)
}

func TestNewCmdList_RunE(t *testing.T) {
	f, _, _ := factory.NewTestFactory()

	var sawOpts *ListOptions
	cmd := NewCmdList(f, func(opts *ListOptions) error {
		sawOpts = opts
		return nil
	})

	cmd.SetArgs([]string{"--limit", "5", "--sort", "desc"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if sawOpts == nil {
		t.Fatal("runF was never called")
	}
	if sawOpts.Limit != 5 {
		t.Errorf("Limit = %d; want 5", sawOpts.Limit)
	}
	if sawOpts.Sort != "desc" {
		t.Errorf("Sort = %q; want %q", sawOpts.Sort, "desc")
	}
}

func TestRunList_ConfigError(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &ListOptions{
		IO:     io,
		Config: func() (*config.Config, error) { return nil, errors.New("no config") },
		Client: func() (*asana.Client, error) { return nil, nil },
	}
	if err := runList(opts); err == nil || !strings.Contains(err.Error(), "failed to get config") {
		t.Fatalf("expected config error, got %v", err)
	}
}

func TestRunList_ClientError(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &ListOptions{
		IO:     io,
		Config: func() (*config.Config, error) { return &config.Config{Workspace: &asana.Workspace{ID: "W"}}, nil },
		Client: func() (*asana.Client, error) { return nil, errors.New("auth failed") },
	}
	if err := runList(opts); err == nil || !strings.Contains(err.Error(), "auth failed") {
		t.Fatalf("expected client error, got %v", err)
	}
}

func TestRunList_IncludesUserIDsInOutput(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	mock, err := asana.NewMockClient(http.StatusOK, []*asana.User{
		{
			ID:   "1198885167962969",
			Name: "Alex Damache",
		},
		{
			ID:   "1198885167962970",
			Name: "Jane Doe",
		},
	})
	if err != nil {
		t.Fatalf("NewMockClient: %v", err)
	}

	opts := &ListOptions{
		IO: io,
		Config: func() (*config.Config, error) {
			return &config.Config{
				Workspace: &asana.Workspace{
					ID:   "W123",
					Name: "TestWS",
				},
			}, nil
		},
		Client: func() (*asana.Client, error) { return newTestClient(mock), nil },
	}

	if err := runList(opts); err != nil {
		t.Fatalf("runList: %v", err)
	}

	got := stripANSI(out.String())
	for _, want := range []string{
		"Alex Damache (id: 1198885167962969/1198885167962969)",
		"Jane Doe (id: 1198885167962970/1198885167962970)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}

func stripANSI(s string) string {
	ansiPattern := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiPattern.ReplaceAllString(s, "")
}
