package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jpoz/goes/cmd/goes/generate"
)

func main() {
	code := run(os.Stdout, os.Args)
	if code != 0 {
		os.Exit(code)
	}
}

const usageText = `usage: goes <command> [<args>...]

goes - Quickly generate Go code to host files build by esbuild

commands:
  generate   Generates Go code from goes.json files
`

func run(w io.Writer, args []string) (code int) {
	if len(args) < 2 {
		fmt.Fprint(w, usageText)
		return 0
	}
	switch args[1] {
	case "generate":
		return generateCmd(w, args[2:])
	}
	fmt.Fprint(w, usageText)

	return 0
}

const generateUsageText = `usage: goes generate [<args>...]

Generates Go code from templ files.

Args:
  -path <path>
    Generates code for all files in path. (default .)
  -help
    Print help and exit.

Examples:

  Generate code for all files in the current directory and subdirectories:

    goes generate

  Generate code in the pkg directory:

    goes generate -path pkg
`

func generateCmd(w io.Writer, args []string) (code int) {
	cmd := flag.NewFlagSet("generate", flag.ExitOnError)
	cmd.SetOutput(w)
	pathFlag := cmd.String("path", ".", "")
	helpFlag := cmd.Bool("help", false, "")

	err := cmd.Parse(args)
	if err != nil || *helpFlag {
		fmt.Fprint(w, generateUsageText)
		return
	}

	err = generate.Run(w, generate.Arguments{
		Path: *pathFlag,
	})
	if err != nil {
		fmt.Fprintln(w, err)
		return 1
	}

	return 0
}
