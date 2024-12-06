package render

import (
	"fmt"
	"io"
)

// PlainFormat implements the Formatter interface to format output in a plain
// text format.
type PlainFormat struct {
	first bool
}

// NewPlainFormat returns a new PlainFormat instance.
func NewPlainFormat() *PlainFormat {
	return &PlainFormat{
		first: true,
	}
}

func (*PlainFormat) WriteHeader(_ io.Writer) error {
	return nil
}

func (p *PlainFormat) WriteBody(w io.Writer, path string, content []byte) error {
	if !p.first {
		fmt.Fprint(w, "----\n")
	}

	p.first = false

	_, err := fmt.Fprintf(w, "%s\n%s\n", path, content)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (*PlainFormat) WriteFooter(_ io.Writer) error {
	return nil
}
