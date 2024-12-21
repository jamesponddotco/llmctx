// Package fscan provides functionality for scanning directory structures
// and collecting files based on configurable rules.
package fscan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~jamesponddotco/gitignore-go"
	"git.sr.ht/~jamesponddotco/llmctx/internal/document"
)

// Directory represents a directory scanning operation with its configuration.
type Directory struct {
	// Root is the root directory to start scanning from.
	Root string

	// GitIgnore is the gitignore matcher to use for filtering files.
	GitIgnore *gitignore.File

	// IgnorePatterns is a list of patterns to ignore during scanning.
	IgnorePatterns []string

	// ShowHidden determines whether hidden files and directories should be
	// included.
	ShowHidden bool
}

// Scan performs the directory scanning according to the configuration
// and returns the collected documents.
func (d *Directory) Scan() (*document.Collection, error) {
	collection := document.New()

	err := filepath.Walk(d.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %s: %w", path, err)
		}

		if path == d.Root {
			return nil
		}

		relPath, err := filepath.Rel(d.Root, path)
		if err != nil {
			return fmt.Errorf("error finding relative path: %w", err)
		}

		skip, err := d.shouldSkip(relPath)
		if err != nil {
			return err
		}

		if skip {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if !info.IsDir() {
			var content []byte

			content, err = os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}

			collection.Add(relPath, content)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return collection, nil
}

// shouldSkip determines if a file or directory should be skipped based on the
// configured rules.
func (d *Directory) shouldSkip(relPath string) (bool, error) {
	if !d.ShowHidden {
		parts := strings.Split(relPath, string(filepath.Separator))
		for _, part := range parts {
			if strings.HasPrefix(part, ".") {
				return true, nil
			}
		}
	}

	for _, pattern := range d.IgnorePatterns {
		matched, err := filepath.Match(pattern, relPath)
		if err != nil {
			return false, fmt.Errorf("error matching pattern %s: %w", pattern, err)
		}

		if matched {
			return true, nil
		}
	}

	if d.GitIgnore != nil {
		if d.GitIgnore.Match(relPath) {
			return true, nil
		}
	}

	return false, nil
}
