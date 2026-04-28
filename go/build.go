package main

import (
	"dagger/go/internal/dagger"
)

// Compile executable
func (m *Go) Compile(
	source *dagger.Directory,
	// +default=""
	path string,
) *dagger.File {
	name := "binary"

	return m.Builder(source).
		WithExec([]string{
			"go", "build", "-o", name, path,
		}).
		File(name)
}
