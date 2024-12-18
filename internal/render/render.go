// Package render contains functions and utilities for rendering output.
package render

import (
	"fmt"
	"io"

	"git.sr.ht/~jamesponddotco/llmctx/internal/document"
)

// Formatter defines an interface for different output formats.
type Formatter interface {
	// WriteHeader writes the output header.
	WriteHeader(w io.Writer) error

	// WriteBody writes the output body.
	WriteBody(w io.Writer, path string, content []byte) error

	// WriteFooter writes the output footer.
	WriteFooter(w io.Writer) error
}

// WriteOutput writes the output to the specified writer using the specified
// format.
func WriteOutput(w io.Writer, collection *document.Collection, format Formatter) error {
	if err := format.WriteHeader(w); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, doc := range collection.Documents {
		if err := format.WriteBody(w, doc.Path, doc.Content); err != nil {
			return fmt.Errorf("failed to write body for %s: %w", doc.Path, err)
		}
	}

	if err := format.WriteFooter(w); err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	return nil
}
