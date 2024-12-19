// Package app is the main package for the application.
package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"git.sr.ht/~jamesponddotco/gitignore-go"
	"git.sr.ht/~jamesponddotco/llmctx/internal/fscan"
	"git.sr.ht/~jamesponddotco/llmctx/internal/meta"
	"git.sr.ht/~jamesponddotco/llmctx/internal/render"
	"git.sr.ht/~jamesponddotco/xstd-go/xflag"
)

// Usage returns the usage information for the application.
func Usage(w io.Writer) {
	text := `NAME:
   %s - %s

USAGE:
   %s [global options]

VERSION:
   %s

GLOBAL OPTIONS:
   --input value, -i value    the directory path to convert (defaults to current directory)
   --output value, -o value   the output file path (defaults to stdout)
   --claude, -c               output in Claude's XML format (defaults to false)
   --show-hidden, -a          show hidden files and directories (defaults to false)
   --ignore value, -x value   patterns to ignore (can be repeated)
   --ignore-gitignore, -g     ignore .gitignore rules (defaults to false)
   --help, -h                 show help
   --version, -v              print the version
`

	fmt.Fprintf(w, text, meta.Name, meta.Description, meta.Name, meta.Version)
}

// Run is the entry point for the application.
func Run(args []string) int {
	var (
		ignoreGitignore bool
		showHidden      bool
		ignorePatterns  xflag.StringSlice
		input           string
		output          string
		claude          bool
		help            bool
		version         bool
	)

	flags := flag.NewFlagSet(meta.Name, flag.ExitOnError)
	flags.StringVar(&input, "input", ".", "the directory path to convert")
	flags.StringVar(&input, "i", ".", "the directory path to convert")
	flags.StringVar(&output, "output", "", "the output txt file path")
	flags.StringVar(&output, "o", "", "the output txt file path")
	flags.BoolVar(&claude, "claude", false, "output in Claude's XML format")
	flags.BoolVar(&claude, "c", false, "output in Claude's XML format")
	flags.BoolVar(&ignoreGitignore, "ignore-gitignore", false, "ignore .gitignore rules")
	flags.BoolVar(&ignoreGitignore, "g", false, "ignore .gitignore rules")
	flags.BoolVar(&showHidden, "show-hidden", false, "show hidden files and directories")
	flags.BoolVar(&showHidden, "a", false, "show hidden files and directories")
	flags.Var(&ignorePatterns, "ignore", "patterns to ignore (can be repeated)")
	flags.Var(&ignorePatterns, "x", "patterns to ignore (can be repeated)")
	flags.BoolVar(&help, "help", false, "show help information")
	flags.BoolVar(&help, "h", false, "show help information")
	flags.BoolVar(&version, "version", false, "print the version")
	flags.BoolVar(&version, "v", false, "print the version")

	flags.Usage = func() {
		Usage(os.Stderr)
	}

	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)

		return 1
	}

	if help {
		Usage(os.Stdout)

		return 0
	}

	if version {
		fmt.Fprintf(os.Stdout, "%s\n", meta.Version)

		return 0
	}

	var out io.Writer

	if output != "" {
		file, err := os.Create(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)

			return 1
		}

		defer func() {
			if err = file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}
		}()

		out = file
	} else {
		out = os.Stdout
	}

	var format render.Formatter
	if claude {
		format = render.NewClaudeFormat()
	} else {
		format = render.NewPlainFormat()
	}

	var (
		matcher       *gitignore.File
		gitignorePath = filepath.Join(input, ".gitignore")
	)

	if _, err := os.Stat(gitignorePath); err == nil && !ignoreGitignore {
		matcher, err = gitignore.New(gitignorePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)

			return 1
		}
	}

	dir := fscan.Directory{
		Root:           input,
		GitIgnore:      matcher,
		IgnorePatterns: ignorePatterns,
		ShowHidden:     showHidden,
	}

	collection, err := dir.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)

		return 1
	}

	if err = render.WriteOutput(out, collection, format); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)

		return 1
	}

	return 0
}
