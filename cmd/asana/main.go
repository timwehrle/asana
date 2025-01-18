package main

import (
	"github.com/timwehrle/asana/pkg/cmd"
	"os"
)

func main() {
	code := cmd.Main()
	os.Exit(int(code))
}
