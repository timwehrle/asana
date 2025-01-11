package factory

import (
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type Factory interface {
	Config() (*config.Config, error)
	NewAsanaClient() (*asana.Client, error)
	Prompter() prompter.Prompter
	IOStreams() *iostreams.IOStreams
}

type DefaultFactory struct {
	prompter prompter.Prompter
	io       *iostreams.IOStreams
}

func New() *DefaultFactory {
	return &DefaultFactory{
		prompter: prompter.New(),
		io:       iostreams.System(),
	}
}

func (f *DefaultFactory) Config() (*config.Config, error) {
	cfg := &config.Config{}

	if err := cfg.Load(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (f *DefaultFactory) NewAsanaClient() (*asana.Client, error) {
	token, err := auth.Get()
	if err != nil {
		return nil, err
	}

	return asana.NewClientWithAccessToken(token), nil
}

func (f *DefaultFactory) Prompter() prompter.Prompter {
	return f.prompter
}

func (f *DefaultFactory) IOStreams() *iostreams.IOStreams {
	return f.io
}
