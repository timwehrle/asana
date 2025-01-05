package main

import (
	"fmt"
	"github.com/timwehrle/asana/cmd"
	"os"
)

func main() {
	root, err := cmd.NewCmdRoot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		fmt.Printf("%s\n\n", err)
		_ = root.Usage()
	}
}
