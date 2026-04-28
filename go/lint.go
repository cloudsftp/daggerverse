package main

import (
	"context"

	"dagger/go/internal/dagger"
)

func (m *Go) linter(
	source *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		From("golangci/golangci-lint:v"+m.GolangCiVersion+"-alpine").

		// Sources
		WithMountedDirectory("/src", source).
		WithWorkdir("/src")
}

func (m *Go) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +default="./..."
	path string,
) error {
	_, err := m.linter(source).
		WithExec([]string{
			"golangci-lint", "run",
			resolvePath(path),
		}).
		Sync(ctx)

	return err
}
