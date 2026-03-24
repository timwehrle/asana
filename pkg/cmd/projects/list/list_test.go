package list

import (
	"errors"
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

	cmd.SetArgs([]string{"--limit", "5", "--favorite"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if sawOpts == nil {
		t.Fatal("runF was never called")
	}
	if sawOpts.Limit != 5 {
		t.Errorf("Limit = %d; want 5", sawOpts.Limit)
	}
	if !sawOpts.Favorite {
		t.Error("Favorite = false; want true")
	}
}

func TestNewCmdList_RunE_InvalidLimit(t *testing.T) {
	f, _, _ := factory.NewTestFactory()
	cmd := NewCmdList(f, func(opts *ListOptions) error { return nil })
	cmd.SetArgs([]string{"--limit", "-1"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "invalid limit") {
		t.Fatalf("expected invalid-limit error, got %v", err)
	}
}

func TestRunList_ConfigError(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &ListOptions{
		IO:     io,
		Config: func() (*config.Config, error) { return nil, errors.New("no config") },
		Client: func() (*asana.Client, error) { return nil, nil },
	}
	if err := runList(opts); err == nil {
		t.Fatal("expected config error, got nil")
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

func TestRunList_IncludesProjectIDsInOutput(t *testing.T) {
	io, _, out, _ := iostreams.Test()

	mock, err := asana.NewMockClient(http.StatusOK, []*asana.Project{
		{
			ID: "1213736174197282",
			ProjectBase: asana.ProjectBase{
				Name: "Agents",
			},
		},
		{
			ID: "1213736174197283",
			ProjectBase: asana.ProjectBase{
				Name: "Backend",
			},
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

	got := out.String()
	for _, want := range []string{
		"Agents (id: 1213736174197282/1213736174197282)",
		"Backend (id: 1213736174197283/1213736174197283)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}
