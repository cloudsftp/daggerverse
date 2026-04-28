package main

import (
	"context"

	"dagger/go/internal/dagger"
)

// Build a go service image
func (m *Go) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
	// +default=""
	path string,
	// +default="program"
	pkg string,
) *dagger.Container {
	executable := m.Compile(source, path)

	a := dag.Alpine()
	return a.ServiceContainer(executable, dagger.AlpineServiceContainerOpts{
		Name: pkg,
	})
}
