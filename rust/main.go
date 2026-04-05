package main

import (
	"fmt"

	"dagger/rust-runner/internal/dagger"
)

type Rust struct {
	RustVersion   string // +default="1.94"
	AlpineVersion string // +default="3.23"
}

// Returns a cached Rust builder container
func (m *Rust) Builder(source *dagger.Directory) *dagger.Container {
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

// Build a service executable
func (m *Rust) BuildExecutable(source *dagger.Directory, name string) *dagger.File {
	return m.Builder(source).
		WithExec([]string{"cargo", "build", "--release", "-p", name}).
		WithExec([]string{"cp", "target/release/" + name, "/" + name}).
		File("/" + name)
}

// Build a service image
func (m *Rust) BuildImage(source *dagger.Directory, name string) *dagger.Container {
	return m.ServiceContainer(m.BuildExecutable(source, name)).
		WithEntrypoint([]string{"/server"})
}

// Create a minimal service container from an executable
func (m *Rust) ServiceContainer(executable *dagger.File) *dagger.Container {
	return dag.Container().
		From("alpine:"+m.AlpineVersion).
		WithFile("/server", executable)
}
