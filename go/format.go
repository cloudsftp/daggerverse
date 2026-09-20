package main

import (
	"context"
	"fmt"

	"dagger/go/internal/dagger"
)

// Check the format of the go code
func (m *Go) Fmt(
	ctx context.Context,
	source *dagger.Directory,
	// +default="./..."
	path string,
) error {
	formattedSource := m.Builder(source).
		WithExec([]string{"go", "fmt", resolveRecursivePath(path)}).
		Directory(SourcePath)

	changes := formattedSource.Changes(source)
	noChanges, err := changes.IsEmpty(ctx)
	if err != nil {
		return fmt.Errorf("could not check changes after formatting go code: %w", err)
	}

	if !noChanges {
		modifiedFiles, err := changes.ModifiedPaths(ctx)
		if err != nil {
			return fmt.Errorf("could not get modified files paths: %w", err)
		}

		return fmt.Errorf("files not formatted correctly: %v", modifiedFiles)
	}

	return nil
}
