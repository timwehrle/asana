package prompter

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultPrompter_Input(t *testing.T) {
	originalAsk := ask

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAsk := func(q survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				if input, ok := response.(*string); ok {
					*input = tt.mockResponse
					return nil
				}
				return errors.New("invalid response type")
			}

			ask = mockAsk
			t.Cleanup(func() {
				ask = originalAsk
			})

			prompter := New()
			result, err := prompter.Input(tt.prompt, tt.defaultValue)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, result)
			}
		})
	}
}
