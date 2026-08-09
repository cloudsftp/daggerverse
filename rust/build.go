package main

import (
	"dagger/rust/internal/dagger"
)

const resultDir = "/result"

// Build a service executable with zig
func (m *Rust) BuildExecutableZig(
	source *dagger.Directory,
	pkg string,
	// +optional
	target string,
) *dagger.File {
	buildCommand := []string{
		"cargo", "zigbuild", "--release",
		"-p", pkg,
	}

	if target != "" {
		buildCommand = append(buildCommand, "--target", target)
	}

	return m.Builder(source).
		WithExec([]string{"cargo", "install", "cargo-zigbuild"}).
		WithExec(buildCommand).
		WithDirectory(resultDir, dag.Directory()).
		WithExec([]string{"cp", "target/release/" + pkg, resultDir}).
		File(resultDir + "/" + pkg)
}

// Build a service executable
func (m *Rust) BuildExecutable(
	source *dagger.Directory,
	pkg string,
	// +optional
	target string,
) *dagger.File {
	buildCommand := []string{
		"cargo", "build", "--release",
		"-p", pkg,
	}

	if target != "" {
		buildCommand = append(buildCommand, "--target", target)
	}

	return m.Builder(source).
		WithExec(buildCommand).
		WithDirectory(resultDir, dag.Directory()).
		WithExec([]string{"cp", "target/release/" + pkg, resultDir}).
		File(resultDir + "/" + pkg)
}
