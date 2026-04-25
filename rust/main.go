package main

import (
	"fmt"

	"dagger/rust/internal/dagger"
)

type Rust struct {
	RustVersion   string
	AlpineVersion string
}

func New(
	// +default="1.95"
	rustVersion string,
	// +default="3.23"
	alpineVersion string,
) *Rust {
	return &Rust{
		RustVersion:   rustVersion,
		AlpineVersion: alpineVersion,
	}
}

// Returns a cached Rust builder container
func (m *Rust) builder(source *dagger.Directory) *dagger.Container {
	source = source.WithoutDirectory("target")

	return dag.Container().
		From(fmt.Sprintf("rust:%s-alpine%s", m.RustVersion, m.AlpineVersion)).
		WithExec([]string{
			"apk", "add", "--no-cache",
			"pkgconfig", "musl-dev",
			"openssl-dev", "openssl-libs-static",
		}).
		WithExec([]string{"rustup", "component", "add", "clippy"}).
		WithDirectory("/src", source).
		WithWorkdir("/src").

		// Caches
		WithMountedCache("/cache/cargo", dag.CacheVolume("rust-packages")).
		WithEnvVariable("CARGO_HOME", "/cache/cargo").
		WithMountedCache("/src/target", dag.CacheVolume("rust-target"))
}
