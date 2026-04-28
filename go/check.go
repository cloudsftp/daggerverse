package main

import (
	"context"

	"dagger/go/internal/dagger"
)

// Check the rust code
func (m *Go) Vet(
	ctx context.Context,
	source *dagger.Directory,
	// +default="./..."
	path string,
) error {
	_, err := m.Builder(source).
		WithExec([]string{
			"go", "vet",
			resolvePath(path),
		}).
		Sync(ctx)

	return err
}
