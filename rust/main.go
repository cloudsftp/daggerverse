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
	// Rust targets to install
	Targets []string
	// Rust crates to install in addition to cargo-zigbuild
	Crates []string
	// Alpine packages to install in addition to pkgconfig and musl-dev
	Packages []string
}

func New(
	// +default="1.97"
	rustVersion string,
	// +default="3.24"
	alpineVersion string,
	// +optional
	components []string,
	// +optional
	targets []string,
	// +optional
	crates []string,
	// +optional
	packages []string,
) *Rust {
	components = append(components, "rustfmt", "clippy")
	packages = append(packages, "pkgconfig", "musl-dev")
	crates = append(crates, "cargo-zigbuild")

	return &Rust{
		rustVersion,
		alpineVersion,
		components,
		targets,
		crates,
		packages,
	}
}

// Returns a bare rust builder container
func (m *Rust) Container() *dagger.Container {
	c := dag.Container().
		From(fmt.Sprintf("rust:%s-alpine%s", m.RustVersion, m.AlpineVersion))

	if len(m.Packages) > 0 {
		c = c.WithExec(append(
			[]string{"apk", "add", "--no-cache"},
			m.Packages...,
		))
	}

	if len(m.Components) > 0 {
		c = c.WithExec(append(
			[]string{"rustup", "component", "add"},
			m.Components...,
		))
	}

	if len(m.Targets) > 0 {
		c = c.WithExec(append(
			[]string{"rustup", "target", "add"},
			m.Targets...,
		))
	}

	if len(m.Crates) > 0 {
		c = c.WithExec(append(
			[]string{"cargo", "install"},
			m.Crates...,
		))
	}

	return c
}

// Returns a cached rust builder container
func (m *Rust) Builder(source *dagger.Directory) *dagger.Container {
	source = source.WithoutDirectory("target")

	return m.Container().
		// Source
		WithDirectory("/src", source).
		WithWorkdir("/src").

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("/src/target", dag.CacheVolume("rust-target"))
}
