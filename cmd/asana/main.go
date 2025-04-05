package main

import (
	"os"

	"github.com/timwehrle/asana/pkg/cmd"
)

func main() {
	code := cmd.Main()
	os.Exit(int(code))
}
