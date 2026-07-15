package main

import (
	"fmt"

	"dagger/bun/internal/dagger"
)

type Bun struct {
	// Version of bun to use
	BunVersion string
	// Alpine packages to install
	Packages []string
}

func New(
	// +default="1.3"
	bunVersion string,
	// +optional
	packages []string,
) *Bun {
	return &Bun{
		bunVersion,
		packages,
	}
}

// Returns a bare bun builder container
func (m *Bun) Container() *dagger.Container {
	return dag.Container().From(fmt.Sprintf("oven/bun:%s-alpine", m.BunVersion))
}

// Returns a bun builder container
func (m *Bun) Builder(source *dagger.Directory) *dagger.Container {
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
		WithWorkdir("/src")
}
