package main

import (
	"fmt"
	"github.com/timwehrle/asana/cmd"
	"github.com/timwehrle/asana/pkg/factory"
	"os"
)

func main() {
	cmdFactory := factory.New()
	stderr := cmdFactory.IOStreams.ErrOut

	root, err := cmd.NewCmdRoot(*cmdFactory)
	if err != nil {
		fmt.Fprintf(stderr, "failed to create root command: %s\n", err)
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		fmt.Printf("%s\n\n", err)
		_ = root.Usage()
	}
}
