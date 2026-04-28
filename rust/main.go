package main

import (
	"fmt"

	"dagger/rust/internal/dagger"
)

type Rust struct {
	// Version of rust to use
	RustVersion string
	// Version of alpine to use
	AlpineVersion string
	// Rust components to install in addition to rustfmt and clippy
	Components []string
	// Alpine packages to install in addition to pkgconfig and musl-dev
	Packages []string
}

func New(
	// +default="1.95"
	rustVersion string,
	// +default="3.23"
	alpineVersion string,
	// +optional
	components []string,
	// +optional
	packages []string,
) *Rust {
	components = append(components, "rustfmt", "clippy")
	packages = append(packages, "pkgconfig", "musl-dev")

	return &Rust{
		rustVersion,
		alpineVersion,
		components,
		packages,
	}
}

// Returns a cached Rust builder container
func (m *Rust) Builder(source *dagger.Directory) *dagger.Container {
	source = source.WithoutDirectory("target")

	builder := dag.Container().From(fmt.Sprintf("rust:%s-alpine%s", m.RustVersion, m.AlpineVersion))

	if len(m.Packages) > 0 {
		builder = builder.WithExec(append(
			[]string{"apk", "add", "--no-cache"},
			m.Packages...,
		))
	}

	if len(m.Components) > 0 {
		builder = builder.WithExec(append(
			[]string{"rustup", "component", "add"},
			m.Components...,
		))

	}

	builder = builder.
		// Source
		WithDirectory("/src", source).
		WithWorkdir("/src").

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("/src/target", dag.CacheVolume("rust-target"))

	return builder
}
