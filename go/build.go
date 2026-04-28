package main

import (
	"dagger/go/internal/dagger"
)

// Compile executable
func (m *Go) Compile(
	source *dagger.Directory,
	name string,
	// +default=""
	path string,
) *dagger.File {
	return m.Builder(source).
		WithExec([]string{
			"go", "build",
			"-o", name,
			resolvePath(path),
		}).
		File(name)
}
