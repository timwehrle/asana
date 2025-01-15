package prompter

import "github.com/stretchr/testify/mock"

type MockPrompter struct {
	mock.Mock
}

func NewMockPrompter() *MockPrompter {
	return &MockPrompter{}
}

func (m *MockPrompter) Input(prompt, defaultValue string) (string, error) {
	args := m.Called(prompt, defaultValue)
	return args.String(0), args.Error(1)
}

func (m *MockPrompter) Confirm(prompt, defaultValue string) (bool, error) {
	args := m.Called(prompt, defaultValue)
	return args.Bool(0), args.Error(1)
}

func (m *MockPrompter) Token() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockPrompter) Select(message string, options []string) (int, error) {
	args := m.Called(message, options)
	return args.Int(0), args.Error(1)
}

func (m *MockPrompter) Editor(prompt, existingDescription string) (string, error) {
	args := m.Called(prompt, existingDescription)
	return args.String(0), args.Error(1)
}
