package render

import (
	"encoding/xml"
	"fmt"
	"io"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
)

// ClaudeFormat implements the Formatter interface to format output in a format
// expected by the LLM Claude.
type ClaudeFormat struct {
	index int
}

// NewClaudeFormat returns a new ClaudeFormat instance.
func NewClaudeFormat() *ClaudeFormat {
	return &ClaudeFormat{
		index: 1,
	}
}

func (*ClaudeFormat) WriteHeader(w io.Writer) error {
	_, err := fmt.Fprintln(w, "<documents>")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *ClaudeFormat) WriteBody(w io.Writer, path string, content []byte) error {
	doc := struct {
		XMLName xml.Name `xml:"document"`
		Source  string   `xml:"source"`
		Content string   `xml:",innerxml"`
		Index   int      `xml:"index,attr"`
	}{
		Source:  path,
		Content: fmt.Sprintf("\n  <document_content>\n%s  </document_content>", xunsafe.BytesToString(content)),
		Index:   c.index,
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")

	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("%w", err)
	}

	c.index++

	fmt.Fprintln(w)

	return nil
}

func (*ClaudeFormat) WriteFooter(w io.Writer) error {
	_, err := fmt.Fprintln(w, "</documents>")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
