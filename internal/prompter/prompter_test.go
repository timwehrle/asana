package prompter

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func runPrompterTests[T any](
	t *testing.T,
	tests []struct {
		name             string
		mockResponse     T
		prompt           string
		defaultValue     string
		expectErr        bool
		expectedResponse T
	},
	mockFunc func(*T, T) error,
	runTest func(prompter Prompter, prompt, defaultValue string) (T, error),
) {
	originalAsk := ask
	defer func() { ask = originalAsk }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ask = func(q survey.Prompt, response any, opts ...survey.AskOpt) error {
				r, ok := response.(*T)
				if !ok {
					return errors.New("invalid response type")
				}
				return mockFunc(r, tt.mockResponse)
			}

			prompter := New()
			result, err := runTest(prompter, tt.prompt, tt.defaultValue)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, result)
			}
		})
	}
}

func TestDefaultPrompter_Input(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     string
		prompt           string
		defaultValue     string
		expectErr        bool
		expectedResponse string
	}{
		{
			name:             "successful input",
			mockResponse:     "mocked response",
			prompt:           "Enter something:",
			defaultValue:     "default",
			expectErr:        false,
			expectedResponse: "mocked response",
		},
		{
			name:             "empty input falls back to default",
			mockResponse:     "",
			prompt:           "Enter something:",
			defaultValue:     "default",
			expectErr:        false,
			expectedResponse: "default",
		},
		{
			name:             "whitespace input falls back to default",
			mockResponse:     "  ",
			prompt:           "Enter something:",
			defaultValue:     "default",
			expectErr:        false,
			expectedResponse: "default",
		},
	}

	runPrompterTests(
		t,
		tests,
		func(response *string, mockResponse string) error {
			*response = mockResponse
			return nil
		},
		func(prompter Prompter, prompt, defaultValue string) (string, error) {
			return prompter.Input(prompt, defaultValue)
		},
	)
}

func TestDefaultPrompter_Confirm(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     bool
		prompt           string
		defaultValue     string
		expectErr        bool
		expectedResponse bool
	}{
		{
			name:             "successful confirmation",
			mockResponse:     true,
			prompt:           "Are you sure?",
			defaultValue:     "No",
			expectErr:        false,
			expectedResponse: true,
		},
		{
			name:             "unsuccessful confirmation",
			mockResponse:     false,
			prompt:           "Are you sure?",
			defaultValue:     "No",
			expectErr:        false,
			expectedResponse: false,
		},
	}

	runPrompterTests(
		t,
		tests,
		func(response *bool, mockResponse bool) error {
			*response = mockResponse
			return nil
		},
		func(prompter Prompter, prompt, defaultValue string) (bool, error) {
			return prompter.Confirm(prompt, defaultValue)
		},
	)
}

func TestDefaultPrompter_Token(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     string
		prompt           string
		defaultValue     string
		expectErr        bool
		expectedResponse string
	}{
		{
			name:             "successful token input",
			mockResponse:     "mocked response",
			prompt:           "Enter your Personal Access Token:",
			defaultValue:     "",
			expectErr:        false,
			expectedResponse: "mocked response",
		},
	}

	runPrompterTests(
		t,
		tests,
		func(response *string, mockResponse string) error {
			*response = mockResponse
			return nil
		},
		func(prompter Prompter, prompt, defaultValue string) (string, error) {
			return prompter.Token()
		},
	)
}

func TestDefaultPrompter_Select(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     int
		prompt           string
		defaultValue     string
		expectErr        bool
		expectedResponse int
	}{
		{
			name:             "successful selection",
			mockResponse:     1,
			prompt:           "Select an option:",
			defaultValue:     "",
			expectErr:        false,
			expectedResponse: 1,
		},
		{
			name:             "successful selection with default value",
			mockResponse:     1,
			prompt:           "Select an option:",
			defaultValue:     "Option 2",
			expectErr:        false,
			expectedResponse: 1,
		},
	}

	runPrompterTests(
		t,
		tests,
		func(response *int, mockResponse int) error {
			*response = mockResponse
			return nil
		},
		func(prompter Prompter, prompt, defaultValue string) (int, error) {
			return prompter.Select(prompt, []string{"Option 1", "Option 2"})
		},
	)
}
