package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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
	/*
		container := .

			// Linter
			WithExec([]string{
				"go", "install",
				"github.com/golangci/golangci-lint/cmd/golangci-lint@" + golangciVersion,
			})
	*/

	return &Go{
		goVersion,
		alpineVersion,
		golangciVersion,
	}
}

// Returns a cached Go builder container
func (m *Go) Builder(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s-alpine%s", m.GoVersion, m.AlpineVersion)).
		WithExec([]string{
			"apk", "add", "--no-cache", "git",
		}).

		// Caches
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache").

		// Source
		WithDirectory("/src", source).
		WithWorkdir("/src")
}

// Build a service executable
func (m *Go) BuildExecutable(
	source *dagger.Directory,
	path string,
) *dagger.File {
	fileName := filepath.Base(path)
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	return m.Builder(source).
		WithExec([]string{
			"go", "build", "-o", name, path,
		}).
		File(fileName)
}

// Build a service image
func (m *Go) BuildImage(
	ctx context.Context,
	source *dagger.Directory,
	path string,
	name string,
) (*dagger.Container, error) {
	executable := m.BuildExecutable(source, path)
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
