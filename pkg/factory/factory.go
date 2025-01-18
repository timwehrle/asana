package factory

import (
	"github.com/timwehrle/asana-api"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type Factory interface {
	Config() (*config.Config, error)
	Client() (*asana.Client, error)
	Prompter() prompter.Prompter
	IOStreams() *iostreams.IOStreams
}

type Default struct {
	prompter prompter.Prompter
	io       *iostreams.IOStreams
}

func New() *Default {
	return &Default{
		prompter: prompter.New(),
		io:       iostreams.System(),
	}
}

func (d *Default) Config() (*config.Config, error) {
	cfg := &config.Config{}

	if err := cfg.Load(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (d *Default) Client() (*asana.Client, error) {
	token, err := auth.Get()
	if err != nil {
		return nil, err
	}

	return asana.NewClientWithAccessToken(token), nil
}

func (d *Default) Prompter() prompter.Prompter {
	return d.prompter
}

func (d *Default) IOStreams() *iostreams.IOStreams {
	return d.io
}
