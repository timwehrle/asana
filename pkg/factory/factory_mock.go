package factory

import (
	"github.com/stretchr/testify/mock"
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/internal/config"
	"github.com/timwehrle/asana/internal/prompter"
	"github.com/timwehrle/asana/pkg/iostreams"
)

type MockFactory struct {
	mock.Mock
	mockPrompter prompter.Prompter
	mockIO       *iostreams.IOStreams
}

func NewMockFactory() *MockFactory {
	return &MockFactory{
		mockPrompter: &prompter.MockPrompter{},
		mockIO: &iostreams.IOStreams{
			In:     iostreams.System().In,
			Out:    iostreams.System().Out,
			ErrOut: iostreams.System().ErrOut,
		},
	}
}

func (f *MockFactory) Config() (*config.Config, error) {
	args := f.Called()
	if cfg, ok := args.Get(0).(*config.Config); ok {
		return cfg, args.Error(1)
	}
	return nil, args.Error(1)
}

func (f *MockFactory) NewAsanaClient() (*asana.Client, error) {
	args := f.Called()
	if client, ok := args.Get(0).(*asana.Client); ok {
		return client, args.Error(1)
	}
	return nil, args.Error(1)
}

func (f *MockFactory) Prompter() prompter.Prompter {
	if args := f.Called(); args.Get(0) != nil {
		return args.Get(0).(prompter.Prompter)
	}
	return f.mockPrompter
}

func (f *MockFactory) IOStreams() *iostreams.IOStreams {
	if args := f.Called(); args.Get(0) != nil {
		return args.Get(0).(*iostreams.IOStreams)
	}
	return f.mockIO
}

func (f *MockFactory) SetPrompter(p prompter.Prompter) {
	f.mockPrompter = p
}

func (f *MockFactory) SetIOStreams(io *iostreams.IOStreams) {
	f.mockIO = io
}
