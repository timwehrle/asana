package utils

import "github.com/fatih/color"

// Colors

func Red() *color.Color {
	return color.New(color.FgRed)
}

func Yellow() *color.Color {
	return color.New(color.FgYellow)
}

func Green() *color.Color {
	return color.New(color.FgGreen)
}

// Styling

func Bold() *color.Color {
	return color.New(color.Bold)
}

func Underline() *color.Color {
	return color.New(color.Underline)
}

func BoldUnderline() *color.Color {
	return color.New(color.Bold, color.Underline)
}

// Status symbols

func Success() string {
	return Green().Sprintf("✓")
}

func Error() string {
	return Red().Sprintf("✗")
}

func Warning() string {
	return Yellow().Sprintf("!")
}
