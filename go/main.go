package main

import (
	"context"
	"fmt"

	"dagger/go-runner/internal/dagger"
)

type Go struct {
	GoVersion       string
	AlpineVersion   string
	GolangCiVersion string
}

func New(
	// +default="1.26"
	goVersion string,
	// +default="3.23"
	alpineVersion string,
	// +default="2.11"
	golangciVersion string,
) *Go {
	return &Go{
		goVersion,
		alpineVersion,
		golangciVersion,
	}
}

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
	name := "program"

	return m.builder(source).
		WithExec([]string{
			"go", "build", "-o", name, path,
		}).
		File(name)
}

// Build a service image
func (m *Go) BuildImage(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
	// +default=""
	path string,
) (*dagger.Container, error) {
	executable := m.Compile(source, path)
	name, err := executable.Name(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get name of executable: %w", err)
	}

	executablePath := "/" + name
	return dag.Container().
		From("alpine:"+m.AlpineVersion).
		WithFile(executablePath, executable).
		WithEntrypoint([]string{executablePath}), nil
}

func (m *Go) linter(
	source *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		From("golangci/golangci-lint:v"+m.GolangCiVersion+"-alpine").

		// Sources
		WithMountedDirectory("/src", source).
		WithWorkdir("/src")
}

func (m *Go) Lint(
	// +defaultPath="/"
	source *dagger.Directory,
	// +default="./..."
	path string,
) *dagger.Container {
	return m.linter(source).
		WithExec([]string{"golangci-lint", "run", path})
}
