package main

import (
	"os"
	"path/filepath"

	"github.com/marcsauter/igcstat/cmd"
)

func main() {
	os.Chdir(filepath.Dir(os.Args[0]))
	cmd.Execute()
}
