// Package document provides types and utilities for handling file content and
// metadata.
package document

// Content represents a single file's content and metadata.
type Content struct {
	// Path is the relative path to the file
	Path string

	// Content is the file's content
	Content []byte
}

// Collection represents a group of documents collected during traversal.
type Collection struct {
	// Documents contains all the files found during traversal
	Documents []Content
}

// New creates a new empty Collection.
func New() *Collection {
	return &Collection{
		Documents: make([]Content, 0),
	}
}

// Add adds a new document to the collection.
func (c *Collection) Add(path string, content []byte) {
	c.Documents = append(c.Documents, Content{
		Path:    path,
		Content: content,
	})
}
