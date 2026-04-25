package main

import (
	"dagger/rust/internal/dagger"
)

// Build a service executable
func (m *Rust) BuildExecutable(source *dagger.Directory, name string) *dagger.File {
	return m.builder(source).
		WithExec([]string{"cargo", "build", "--release", "-p", name}).
		WithExec([]string{"cp", "target/release/" + name, "/" + name}).
		File("/" + name)
}
