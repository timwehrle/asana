package utils

import (
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	color.NoColor = false
}

func TestColorFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func() *color.Color
	}{
		{"Red", Red},
		{"Yellow", Yellow},
		{"Green", Green},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			formatted := result.Sprint("test")
			assert.NotEqual(t, "test", formatted, "Color formatting should modify the string")
			assert.Contains(t, formatted, "\x1b[", "Should contain ANSI escape sequence")
		})
	}
}

func TestStylingFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func() *color.Color
	}{
		{"Bold", Bold},
		{"Underline", Underline},
		{"BoldUnderline", BoldUnderline},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			formatted := result.Sprint("test")
			assert.NotEqual(t, "test", formatted, "Style formatting should modify the string")
			assert.Contains(t, formatted, "\x1b[", "Should contain ANSI escape sequence")
		})
	}
}

func TestStatusSymbols(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() string
		contains string
	}{
		{"Success", Success, "✓"},
		{"Error", Error, "✗"},
		{"Warning", Warning, "!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			assert.Contains(t, result, tt.contains)
			assert.Contains(t, result, "\x1b[", "Should contain ANSI escape sequence")
		})
	}
}
