package list

import (
	"bytes"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
	"testing"
)

func makeTestFactory() (factory.Factory, *bytes.Buffer, *bytes.Buffer) {
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	ios := &iostreams.IOStreams{In: nil, Out: outBuf, ErrOut: errBuf}
	fakeConfig := func() (*config.Config, error) {
		return &config.Config{
			Workspace: &asana.Workspace{
				ID:   "W123",
				Name: "TestWS",
			},
		}, nil
	}
	fakeClient := func() (*asana.Client, error) {
		return &asana.Client{}, nil
	}
	return factory.Factory{
		IOStreams: ios,
		Config:    fakeConfig,
		Client:    fakeClient,
		Prompter:  nil,
	}, outBuf, errBuf
}

func TestNewCmdList_RunE(t *testing.T) {
	f, _, _ := makeTestFactory()
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
