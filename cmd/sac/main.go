package main

import (
	"os"

	"github.com/s4bb4t/sa-cli/internal/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
