package main

import (
	"fmt"

	"dagger/go/internal/dagger"
)

type Go struct {
	// Version of go to use
	GoVersion string
	// Version of alpine to use
	AlpineVersion string
	// Version of golang ci lint to use
	GolangCiVersion string
	// Alpine packages to install
	Packages []string
}

func New(
	// +default="1.26"
	goVersion string,
	// +default="2.12"
	golangCiVersion string,
	// +default="3.24"
	alpineVersion string,
	// +optional
	packages []string,
) *Go {
	return &Go{
		goVersion,
		alpineVersion,
		golangCiVersion,
		packages,
	}
}

// Returns a bare go builder container
func (m *Go) Container() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s-alpine%s", m.GoVersion, m.AlpineVersion))
}

// Returns a cached go builder container
func (m *Go) Builder(source *dagger.Directory) *dagger.Container {
	builder := m.Container()

	if len(m.Packages) > 0 {
		builder = builder.WithExec(append(
			[]string{"apk", "add", "--no-cache"},
			m.Packages...,
		))
	}

	return builder.
		// Source
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").

		// Caches
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build")).
		WithEnvVariable("GOCACHE", "/go/build-cache")
}
