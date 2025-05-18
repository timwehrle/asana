package factory

import (
	"bytes"
	"github.com/timwehrle/asana/internal/api/asana"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/pkg/iostreams"
)

// NewTestFactory returns a Factory prewired for tests,
// plus buffers for capturing stdout/stderr.
func NewTestFactory() (Factory, *bytes.Buffer, *bytes.Buffer) {
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	io := &iostreams.IOStreams{
		In:     nil,
		Out:    outBuf,
		ErrOut: errBuf,
	}

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

	return Factory{
		IOStreams: io,
		Config:    fakeConfig,
		Client:    fakeClient,
		Prompter:  nil,
	}, outBuf, errBuf
}
