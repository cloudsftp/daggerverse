package main

import (
	"context"

	"dagger/go/internal/dagger"
)

// Build a go service image
func (m *Go) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
	name string,
	// +default=""
	path string,
) *dagger.Container {
	executable := m.Compile(source, name, path)

	a := dag.Alpine()
	return a.ServiceContainer(executable, name)
}
