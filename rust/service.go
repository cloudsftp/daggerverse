package main

import (
	"context"

	"dagger/rust/internal/dagger"
)

// Build a service image
func (m *Rust) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
	pkg string,
) *dagger.Container {
	executable := m.BuildExecutable(source, pkg)

	a := dag.Alpine()
	return a.ServiceContainer(executable, pkg)
}
