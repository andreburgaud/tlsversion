package main

import (
	"os"

	"tlsversion/internal/cli"

	"github.com/fatih/color"
)

func main() {
	cmd, err := cli.ParseOptions()
	red := color.New(color.FgHiRed, color.Bold)
	if err != nil {
		_, _ = red.Println(err)
		os.Exit(1)
	}

	err = cmd.Execute()
	if err != nil {
		_, _ = red.Println(err)
		os.Exit(1)
	}
}
