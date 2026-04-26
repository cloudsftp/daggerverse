package main

import (
	"context"
	"fmt"

	"dagger/rust/internal/dagger"
)

// Build a service image
func (m *Rust) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
	name string,
) (*dagger.Container, error) {
	executable, err := m.BuildExecutable(ctx, source, name)
	if err != nil {
		return nil, fmt.Errorf("could not build executable: %w", err)
	}

	return m.ServiceContainer(executable, name).
		WithEntrypoint([]string{"/" + name}), nil
}

// Create a minimal service container from an executable
func (m *Rust) ServiceContainer(
	executable *dagger.File,
	name string,
) *dagger.Container {
	return dag.Container().
		From("alpine:"+m.AlpineVersion).
		WithFile("/"+name, executable)
}
