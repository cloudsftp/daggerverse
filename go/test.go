package main

import (
	"context"

	"dagger/go/internal/dagger"
)

// Run tests
func (m *Go) Test(
	ctx context.Context,
	source *dagger.Directory,
	// +default="./..."
	path string,
) error {
	_, err := m.Builder(source).
		WithExec([]string{
			"go", "test",
			resolveRecursivePath(path),
		}).
		Sync(ctx)

	return err
}
