package factory

import (
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/internal/config"
)

type Factory interface {
	Config() (*config.Config, error)
	NewAsanaClient() (*asana.Client, error)
}

type DefaultFactory struct{}

func New() *DefaultFactory {
	return &DefaultFactory{}
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
