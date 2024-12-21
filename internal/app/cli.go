// Package app is the main package for the application.
package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"git.sr.ht/~jamesponddotco/errxit-go"
	"git.sr.ht/~jamesponddotco/gitignore-go"
	"git.sr.ht/~jamesponddotco/llmctx/internal/fscan"
	"git.sr.ht/~jamesponddotco/llmctx/internal/meta"
	"git.sr.ht/~jamesponddotco/llmctx/internal/render"
	"git.sr.ht/~jamesponddotco/xstd-go/xflag"
)

// CLI is the command line interface for the application.
type CLI struct {
	// Input is the path to the input directory.
	Input string

	// Output is the path to the output file.
	Output string

	// IgnorePatterns is a list of patterns to ignore when scanning.
	IgnorePatterns []string

	// IgnoreGitignore tells the application to ignore the .gitignore file.
	IgnoreGitignore bool

	// ShowHidden tells the application to show hidden files and directories.
	ShowHidden bool

	// Claude tells the application to output in Claude's XML format.
	Claude bool

	// Help tells the application to show the help message.
	Help bool

	// Version tells the application to show the version.
	Version bool
}

// New parses the command line arguments and returns a new CLI instance.
func New(args []string) CLI {
	var (
		cli            CLI
		ignorePatterns xflag.StringSlice
	)

	flags := flag.NewFlagSet(meta.Name, flag.ExitOnError)
	flags.StringVar(&cli.Input, "input", ".", "the directory path to convert")
	flags.StringVar(&cli.Input, "i", ".", "the directory path to convert")
	flags.StringVar(&cli.Output, "output", "", "the output txt file path")
	flags.StringVar(&cli.Output, "o", "", "the output txt file path")
	flags.BoolVar(&cli.Claude, "claude", false, "output in Claude's XML format")
	flags.BoolVar(&cli.Claude, "c", false, "output in Claude's XML format")
	flags.BoolVar(&cli.IgnoreGitignore, "ignore-gitignore", false, "ignore .gitignore rules")
	flags.BoolVar(&cli.IgnoreGitignore, "g", false, "ignore .gitignore rules")
	flags.BoolVar(&cli.ShowHidden, "show-hidden", false, "show hidden files and directories")
	flags.BoolVar(&cli.ShowHidden, "a", false, "show hidden files and directories")
	flags.Var(&ignorePatterns, "ignore", "patterns to ignore (can be repeated)")
	flags.Var(&ignorePatterns, "x", "patterns to ignore (can be repeated)")
	flags.BoolVar(&cli.Help, "help", false, "show help information")
	flags.BoolVar(&cli.Help, "h", false, "show help information")
	flags.BoolVar(&cli.Version, "version", false, "print the version")
	flags.BoolVar(&cli.Version, "v", false, "print the version")

	flags.Usage = func() {
		Usage(os.Stderr)
	}

	flags.Parse(args) //nolint:errcheck // we can't return the error anyway, as we set flag.ExitOnError

	cli.IgnorePatterns = ignorePatterns

	return cli
}

// Run is the entry point for the application.
func (c *CLI) Run() error {
	if c.Help {
		Usage(os.Stdout)

		return nil
	}

	if c.Version {
		fmt.Fprintf(os.Stdout, "%s\n", meta.Version)

		return nil
	}

	var out io.Writer

	if c.Output != "" {
		file, err := os.Create(c.Output)
		if err != nil {
			return &errxit.Error{
				Err: err,
				No:  1,
			}
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
	if c.Claude {
		format = render.NewClaudeFormat()
	} else {
		format = render.NewPlainFormat()
	}

	var (
		matcher       *gitignore.File
		gitignorePath = filepath.Join(c.Input, ".gitignore")
	)

	if _, err := os.Stat(gitignorePath); err == nil && !c.IgnoreGitignore {
		matcher, err = gitignore.New(gitignorePath)
		if err != nil {
			return &errxit.Error{
				Err: err,
				No:  1,
			}
		}
	}

	dir := fscan.Directory{
		Root:           c.Input,
		GitIgnore:      matcher,
		IgnorePatterns: c.IgnorePatterns,
		ShowHidden:     c.ShowHidden,
	}

	collection, err := dir.Scan()
	if err != nil {
		return &errxit.Error{
			Err: err,
			No:  1,
		}
	}

	if err = render.WriteOutput(out, collection, format); err != nil {
		return &errxit.Error{
			Err: err,
			No:  1,
		}
	}

	return nil
}
