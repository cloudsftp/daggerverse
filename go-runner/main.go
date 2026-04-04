package main

import (
	"fmt"

	"dagger/go-runner/internal/dagger"
)

type GoRunner struct {
	GoVersion       string // +default="1.26"
	AlpineVersion   string // +default="3.23"
	GolangciVersion string // +default="v2.11"
}

// Returns a cached Go builder container
func (m *GoRunner) Builder(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s-alpine%s", m.GoVersion, m.AlpineVersion)).

		// Caches
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").

		// Linter
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@" + m.GolangciVersion}).

		// Source
		WithDirectory("/src", source).
		WithWorkdir("/src")
}

// Build a service executable
func (m *GoRunner) BuildExecutable(source *dagger.Directory, path string, name string) *dagger.File {
	return m.Builder(source).
		WithExec([]string{
			"go", "build", "-o", name,
			fmt.Sprintf("%s/%s.go", path, name),
		}).
		File(name)
}

// Build a service image
func (m *GoRunner) BuildImage(source *dagger.Directory, path string, name string) *dagger.Container {
	return m.ServiceContainer(m.BuildExecutable(source, path, name)).
		WithEntrypoint([]string{"/server"})
}

// Create a minimal service container from an executable
func (m *GoRunner) ServiceContainer(executable *dagger.File) *dagger.Container {
	return dag.Container().
		From("alpine:"+m.AlpineVersion).
		WithFile("/server", executable)
}
