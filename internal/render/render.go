// Package render contains functions and utilities for rendering output.
package render

import "io"

// Formatter defines an interface for different output formats.
type Formatter interface {
	// WriteHeader writes the output header.
	WriteHeader(w io.Writer) error

	// WriteBody writes the output body.
	WriteBody(w io.Writer, path string, content []byte) error

	// WriteFooter writes the output footer.
	WriteFooter(w io.Writer) error
}
