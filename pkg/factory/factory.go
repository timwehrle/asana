package factory

import (
	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type Factory struct {
	Config func() (*config.Config, error)
	Client func() (*asana.Client, error)

	Prompter  prompter.Prompter
	IOStreams *iostreams.IOStreams
}

func New() *Factory {
	f := &Factory{}

	f.IOStreams = ioStreams()
	f.Prompter = newPrompter()
	f.Client = newClientFunc()
	f.Config = newConfigFunc()

	return f
}

func newConfigFunc() func() (*config.Config, error) {
	return func() (*config.Config, error) {
		cfg := &config.Config{}

		if err := cfg.Load(); err != nil {
			return nil, err
		}

		return cfg, nil
	}
}

func newClientFunc() func() (*asana.Client, error) {
	return func() (*asana.Client, error) {
		token, err := auth.Get()
		if err != nil {
			return nil, err
		}

		return asana.NewClientWithAccessToken(token), nil
	}
}

func newPrompter() prompter.Prompter {
	return prompter.New()
}

func ioStreams() *iostreams.IOStreams {
	io := iostreams.System()

	return io
}
