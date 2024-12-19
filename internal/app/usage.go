package app

import (
	"fmt"
	"io"

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
