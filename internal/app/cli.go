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
	// Metadata is the metadata for the application.
	Metadata *xflag.Metadata

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

	cli.Metadata = &xflag.Metadata{
		Name:        meta.Name,
		Description: meta.Description,
		Ver:         meta.Version,
		Options: []xflag.Option{
			{
				Name:        "input",
				Shorthand:   "i",
				Type:        "DIR",
				Description: "the directory path to convert",
				Default:     "current directory",
			},
			{
				Name:        "output",
				Shorthand:   "o",
				Type:        "FILE",
				Description: "the output txt file path",
				Default:     "stdout",
			},
			{
				Name:        "claude",
				Shorthand:   "c",
				Description: "output in Claude's XML format",
				Default:     "false",
			},
			{
				Name:        "ignore-gitignore",
				Shorthand:   "g",
				Description: "ignore .gitignore rules",
				Default:     "false",
			},
			{
				Name:        "show-hidden",
				Shorthand:   "a",
				Description: "show hidden files and directories",
				Default:     "false",
			},
			{
				Name:        "ignore",
				Shorthand:   "x",
				Type:        "PATTERN",
				Description: "patterns to ignore (can be repeated)",
			},
			{
				Name:        "help",
				Shorthand:   "h",
				Description: "show help information",
			},
			{
				Name:        "version",
				Shorthand:   "v",
				Description: "print the version",
			},
		},
	}

	flags := flag.NewFlagSet(meta.Name, flag.ExitOnError)
	flags.StringVar(&cli.Input, cli.Metadata.Options[0].Name, ".", cli.Metadata.Options[0].Description)
	flags.StringVar(&cli.Input, cli.Metadata.Options[0].Shorthand, ".", cli.Metadata.Options[0].Description)
	flags.StringVar(&cli.Output, cli.Metadata.Options[1].Name, "", cli.Metadata.Options[1].Description)
	flags.StringVar(&cli.Output, cli.Metadata.Options[1].Shorthand, "", cli.Metadata.Options[1].Description)
	flags.BoolVar(&cli.Claude, cli.Metadata.Options[2].Name, false, cli.Metadata.Options[2].Description)
	flags.BoolVar(&cli.Claude, cli.Metadata.Options[2].Shorthand, false, cli.Metadata.Options[2].Description)
	flags.BoolVar(&cli.IgnoreGitignore, cli.Metadata.Options[3].Name, false, cli.Metadata.Options[3].Description)
	flags.BoolVar(&cli.IgnoreGitignore, cli.Metadata.Options[3].Shorthand, false, cli.Metadata.Options[3].Description)
	flags.BoolVar(&cli.ShowHidden, cli.Metadata.Options[4].Name, false, cli.Metadata.Options[4].Description)
	flags.BoolVar(&cli.ShowHidden, cli.Metadata.Options[4].Shorthand, false, cli.Metadata.Options[4].Description)
	flags.Var(&ignorePatterns, cli.Metadata.Options[5].Name, cli.Metadata.Options[5].Description)
	flags.Var(&ignorePatterns, cli.Metadata.Options[5].Shorthand, cli.Metadata.Options[5].Description)
	flags.BoolVar(&cli.Help, cli.Metadata.Options[6].Name, false, cli.Metadata.Options[6].Description)
	flags.BoolVar(&cli.Help, cli.Metadata.Options[6].Shorthand, false, cli.Metadata.Options[6].Description)
	flags.BoolVar(&cli.Version, cli.Metadata.Options[7].Name, false, cli.Metadata.Options[7].Description)
	flags.BoolVar(&cli.Version, cli.Metadata.Options[7].Shorthand, false, cli.Metadata.Options[7].Description)

	flags.Usage = func() {
		fmt.Fprint(os.Stderr, cli.Metadata.Usage())
	}

	flags.Parse(args) //nolint:errcheck // we can't return the error anyway, as we set flag.ExitOnError

	cli.IgnorePatterns = ignorePatterns

	return cli
}

// Run is the entry point for the application.
func (c *CLI) Run() error {
	if c.Help {
		fmt.Fprint(os.Stderr, c.Metadata.Usage())

		return nil
	}

	if c.Version {
		fmt.Fprintf(os.Stdout, "%s\n", c.Metadata.Version())

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
