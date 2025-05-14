package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/timwehrle/asana/internal/build"
	"github.com/timwehrle/asana/pkg/cmdutils"
	"strings"

	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/mgutz/ansi"
	"github.com/timwehrle/asana/internal/auth"
	"github.com/timwehrle/asana/pkg/cmd/root"
	"github.com/timwehrle/asana/pkg/factory"
)

type ExitCode int

const (
	exitOK     ExitCode = 0
	exitError  ExitCode = 1
	exitCancel ExitCode = 2
	exitAuth   ExitCode = 4
)

func Main() ExitCode {
	buildVersion := build.Version

	f := factory.New()
	stderr := f.IOStreams.ErrOut
	cs := f.IOStreams.ColorScheme()

	if !f.IOStreams.ColorEnabled {
		surveyCore.DisableColor = true
		f.IOStreams.DisableColor()
		ansi.DisableColors(true)
	} else {
		surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
			switch style {
			case "white":
				return ansi.ColorCode("default")
			default:
				return ansi.ColorCode(style)
			}
		}
	}

	rootCmd, err := root.NewCmdRoot(*f, buildVersion)
	if err != nil {
		fmt.Fprintf(stderr, "failed to create root command: %s\n", err)
		return exitError
	}

	rootCmd.PersistentFlags().Bool("help", false, "Show help for command")
	rootCmd.Flags().BoolP("version", "v", false, "Show asana version")

	if err := rootCmd.Execute(); err != nil {
		if cmdutils.IsUserCancellation(err) {
			if errors.Is(err, terminal.InterruptErr) {
				fmt.Fprintf(stderr, "\n")
			}

			return exitCancel
		}

		var authError *auth.AuthenticationError

		if errors.As(err, &authError) {
			return exitAuth
		}

		label := "Error:"
		if f.IOStreams.ColorEnabled {
			label = cs.Error(label)
		}
		fmt.Fprintf(stderr, "%s %s\n\n", label, err)

		if isFlagOrArgError(err) {
			_ = rootCmd.Usage()
		}

		return exitError
	}

	return exitOK
}

func isFlagOrArgError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "unknown flag") || strings.Contains(msg, "requires") || strings.Contains(msg, "invalid") || strings.Contains(msg, "unknown")
}
