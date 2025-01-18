package cmd

import (
	"errors"
	"fmt"
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
	cmdFactory := factory.New()
	stderr := cmdFactory.IOStreams.ErrOut

	if !cmdFactory.IOStreams.ColorEnabled {
		surveyCore.DisableColor = true
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

	rootCmd, err := root.NewCmdRoot(*cmdFactory)
	if err != nil {
		fmt.Fprintf(stderr, "failed to create root command: %s\n", err)
		return exitError
	}

	if err := rootCmd.Execute(); err != nil {
		var authError *auth.AuthenticationError

		if errors.As(err, &authError) {
			return exitAuth
		}

		fmt.Fprintf(stderr, "%s\n\n", err)
		_ = rootCmd.Usage()

		return exitError
	}

	return exitOK
}
