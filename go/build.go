package main

import (
	"fmt"

	"dagger/go/internal/dagger"
)

func (m *Go) builder(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s-alpine%s", m.GoVersion, m.AlpineVersion)).

		// Caches
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").

		// Sources
		WithMountedDirectory("/src", source).
		WithWorkdir("/src")
}

// Compile executable
func (m *Go) Compile(
	// +defaultPath="/"
	source *dagger.Directory,
	// +default=""
	path string,
) *dagger.File {
	name := "binary"

	return m.builder(source).
		WithExec([]string{
			"go", "build", "-o", name, path,
		}).
		File(name)
}
