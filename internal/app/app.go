// Package app is the main package for the application.
package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"git.sr.ht/~jamesponddotco/llmctx/internal/meta"
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
   --input value, -i value  the directory path to convert (defaults to current directory)
   --output value, -o value the output file path (defaults to stdout)
   --help, -h               show help
   --version, -v            print the version
`

	fmt.Fprintf(w, text, meta.Name, meta.Description, meta.Name, meta.Version)
}

// Run is the entry point for the application.
func Run(args []string) int {
	var (
		input   string
		output  string
		help    bool
		version bool
	)

	flags := flag.NewFlagSet(meta.Name, flag.ExitOnError)
	flags.StringVar(&input, "input", ".", "the directory path to convert")
	flags.StringVar(&input, "i", ".", "the directory path to convert")
	flags.StringVar(&output, "output", "", "the output txt file path")
	flags.StringVar(&output, "o", "", "the output txt file path")
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
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}
		}()

		out = file
	} else {
		out = os.Stdout
	}

	if err := WalkDir(input, out); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)

		return 1
	}

	return 0
}

// WalkDir traverses the given directory and writes its structure and file
// contents to the provided io.Writer.
func WalkDir(rootDir string, out io.Writer) error {
	first := true

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %s: %w", path, err)
		}

		if path == rootDir {
			return nil
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return fmt.Errorf("error finding relative path: %w", err)
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			if first {
				first = false
			} else {
				fmt.Fprint(out, "----\n")
			}

			fmt.Fprintf(out, "%s\n%s\n", relPath, string(content))
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk through directory: %w", err)
	}

	return nil
}
