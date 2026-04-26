package main

import (
	"context"
	"fmt"

	"dagger/rust/internal/dagger"
)

const (
	user = "appuser"
)

// Build a service image
func (m *Rust) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
	pkg string,
) (*dagger.Container, error) {
	executable, err := m.BuildExecutable(ctx, source, pkg)
	if err != nil {
		return nil, fmt.Errorf("could not build executable: %w", err)
	}

	return m.ServiceContainer(executable, pkg).
		WithEntrypoint([]string{"/" + pkg}), nil
}

// Create a minimal service container from an executable
func (m *Rust) ServiceContainer(
	executable *dagger.File,
	name string,
) *dagger.Container {
	return dag.Container().
		From("alpine:"+m.AlpineVersion).
		WithExec([]string{"adduser", user, "-D"}).
		WithUser(user).
		WithFile("/"+name, executable)
}
