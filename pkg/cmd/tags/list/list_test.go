package list

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
)

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

func TestFetchFavoriteTags(t *testing.T) {
	tests := []struct {
		name          string
		mockStatus    int
		mockBody      any
		wantErr       bool
		wantTags      []*asana.Tag
		wantPath      string
		wantQueryVals url.Values
	}{
		{
			name:       "success",
			mockStatus: 200,
			mockBody: []*asana.Tag{
				{
					ID: "T1",
					TagBase: asana.TagBase{
						Name: "TagOne",
					},
				},
				{
					ID: "T2",
					TagBase: asana.TagBase{
						Name: "TagTwo",
					},
				},
			},
			wantErr: false,
			wantTags: []*asana.Tag{
				{
					ID: "T1",
					TagBase: asana.TagBase{
						Name: "TagOne",
					},
				},
				{
					ID: "T2",
					TagBase: asana.TagBase{
						Name: "TagTwo",
					},
				},
			},
			wantPath: "/api/1.0/users/me/favorites",
			wantQueryVals: url.Values{
				"resource_type": []string{"tag"},
				"workspace":     []string{"W123"},
			},
		},
		{
			name:       "server error",
			mockStatus: 500,
			mockBody:   "oops",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := asana.NewMockClient(tt.mockStatus, tt.mockBody)
			if err != nil {
				t.Fatalf("NewMockClient: %v", err)
			}
			client := newTestClient(mock)
			ws := &asana.Workspace{ID: "W123"}

			got, err := fetchFavoriteTags(client, ws)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.wantTags) {
				t.Errorf("tags = %#v; want %#v", got, tt.wantTags)
			}

			last := mock.GetLastRequest()
			if got, want := last.Method(), http.MethodGet; got != want {
				t.Errorf("Method = %q; want %q", got, want)
			}
			if got, want := last.Path(), tt.wantPath; got != want {
				t.Errorf("Path = %q; want %q", got, want)
			}
			for key, vals := range tt.wantQueryVals {
				if qv := last.Query()[key]; !reflect.DeepEqual(qv, vals) {
					t.Errorf("query[%q] = %v; want %v", key, qv, vals)
				}
			}
		})
	}
}

func TestFetchFavoriteTags_ErrorPathWrapped(t *testing.T) {
	mock, _ := asana.NewMockClient(500, "fail")
	client := newTestClient(mock)
	ws := &asana.Workspace{ID: "W500"}

	_, err := fetchFavoriteTags(client, ws)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed fetching favorite tags") {
		t.Errorf("error did not wrap correctly: %v", err)
	}
}
