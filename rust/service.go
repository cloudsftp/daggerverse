package main

import (
	"dagger/rust/internal/dagger"
)

// Build a service image
func (m *Rust) BuildImage(source *dagger.Directory, name string) *dagger.Container {
	return m.ServiceContainer(m.BuildExecutable(source, name)).
		WithEntrypoint([]string{"/server"})
}

// Create a minimal service container from an executable
func (m *Rust) ServiceContainer(executable *dagger.File) *dagger.Container {
	return dag.Container().
		From("alpine:"+m.AlpineVersion).
		WithFile("/server", executable)
}
