package main

import (
	"context"

	"dagger/go/internal/dagger"
)

// Run tests
func (m *Go) Test(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="./..."
	path string,
) error {
	_, err := m.builder(source).
		WithExec([]string{
			"go", "test", path,
		}).
		Sync(ctx)

	return err
}
