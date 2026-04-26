package main

import (
	"context"

	"dagger/rust/internal/dagger"
)

const resultDir = "/result"

// Build a service executable
func (m *Rust) BuildExecutable(
	ctx context.Context,
	source *dagger.Directory,
	name string,
) (*dagger.File, error) {
	buildCommand := []string{
		"cargo", "build", "--release",
		"-p", name,
	}

	return m.builder(source).
		WithExec(buildCommand).
		WithDirectory(resultDir, dag.Directory()).
		WithExec([]string{"cp", "target/release/" + name, resultDir}).
		File(resultDir + "/" + name), nil
}
