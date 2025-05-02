package iostreams

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
)

// ColorScheme holds the configured colors for various types of output
type ColorScheme struct {
	// Commands and UI Elements
	Primary   string
	Secondary string
	Success   func(string) string
	Warning   func(string) string
	Error     func(string) string
	Gray      string

	// Special formatting
	Bold        func(string) string
	Dim         func(string) string
	Italic      func(string) string
	Underline   func(string) string
	SuccessIcon string
	WarningIcon string
	ErrorIcon   string
}

// IOStreams provides access to input and output streams
type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	// Whether the output streams are a terminal
	IsStdinTTY  bool
	IsStdoutTTY bool
	IsStderrTTY bool

	// Whether colors should be used
	ColorEnabled bool

	// The color scheme to use
	colorScheme *ColorScheme
}

// System returns an IOStreams suitable for use in a command
func System() *IOStreams {
	ioSys := &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	// Determine if streams are TTYs
	ioSys.IsStdinTTY = isTerminal(os.Stdin)
	ioSys.IsStdoutTTY = isTerminal(os.Stdout)
	ioSys.IsStderrTTY = isTerminal(os.Stderr)

	// Enable colors if stdout is a terminal
	ioSys.ColorEnabled = ioSys.IsStdoutTTY

	// Initialize color scheme
	ioSys.initColorScheme()

	// If on Windows and output is a terminal, use colorable
	if ioSys.IsStdoutTTY {
		ioSys.Out = colorable.NewColorable(os.Stdout)
	}
	if ioSys.IsStderrTTY {
		ioSys.ErrOut = colorable.NewColorable(os.Stderr)
	}

	return ioSys
}

// Test returns an IOStreams suitable for testing
func Test() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	testIO := &IOStreams{
		In:           io.NopCloser(in),
		Out:          out,
		ErrOut:       errOut,
		IsStdinTTY:   false,
		IsStdoutTTY:  false,
		IsStderrTTY:  false,
		ColorEnabled: false,
	}

	testIO.initColorScheme()

	return testIO, in, out, errOut
}

// ColorScheme returns the color scheme for the streams
func (io *IOStreams) ColorScheme() *ColorScheme {
	return io.colorScheme
}

// initColorScheme initializes the color scheme based on whether colors are enabled
func (io *IOStreams) initColorScheme() {
	useColors := io.ColorEnabled

	io.colorScheme = &ColorScheme{
		Primary:   ansi.ColorCode("blue"),
		Secondary: ansi.ColorCode("cyan"),
		Success:   ansi.ColorFunc("green"),
		Warning:   ansi.ColorFunc("yellow"),
		Error:     ansi.ColorFunc("red"),
		Gray:      ansi.ColorCode("black+h"),

		Bold: func(s string) string {
			if !useColors {
				return s
			}
			return ansi.Color(s, "default+b")
		},

		Dim: func(s string) string {
			if !useColors {
				return s
			}
			return ansi.Color(s, "black+h")
		},

		Italic: func(s string) string {
			if !useColors {
				return s
			}
			return ansi.Color(s, "default+i")
		},

		Underline: func(s string) string {
			if !useColors {
				return s
			}
			return ansi.Color(s, "default+u")
		},

		SuccessIcon: func() string {
			if !useColors {
				return "✓"
			}
			return ansi.Color("✓", "green")
		}(),

		WarningIcon: func() string {
			if !useColors {
				return "!"
			}
			return ansi.Color("!", "yellow")
		}(),

		ErrorIcon: func() string {
			if !useColors {
				return "✕"
			}
			return ansi.Color("✕", "red")
		}(),
	}
}

// isTerminal returns true if the given file is a terminal
func isTerminal(f any) bool {
	if file, ok := f.(*os.File); ok {
		return isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd())
	}
	return false
}

// ForceColor forces the use of colors in output
func (io *IOStreams) ForceColor() {
	io.ColorEnabled = true
	io.initColorScheme()
}

// DisableColor disables the use of colors in output
func (io *IOStreams) DisableColor() {
	io.ColorEnabled = false
	io.initColorScheme()
}

// Color returns a string in the given color if colors are enabled
func (io *IOStreams) Color(s string, color string) string {
	if !io.ColorEnabled {
		return s
	}
	return ansi.Color(s, color)
}

// ColorFromScheme returns a string in the given scheme color if colors are enabled
func (io *IOStreams) ColorFromScheme(s string, color func(string) string) string {
	if !io.ColorEnabled {
		return s
	}
	return color(s)
}

// Printf formats and prints to the configured output
func (io *IOStreams) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(io.Out, format, a...)
}

// Println prints to the configured output followed by a newline
func (io *IOStreams) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(io.Out, a...)
}

// ErrPrintf formats and prints to the configured error output
func (io *IOStreams) ErrPrintf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(io.ErrOut, format, a...)
}

// ErrPrintln prints to the configured error output followed by a newline
func (io *IOStreams) ErrPrintln(a ...any) (n int, err error) {
	return fmt.Fprintln(io.ErrOut, a...)
}
