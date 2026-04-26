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
	pkg string,
) (*dagger.File, error) {
	buildCommand := []string{
		"cargo", "build", "--release",
		"-p", pkg,
	}

	return m.builder(source).
		WithExec(buildCommand).
		WithDirectory(resultDir, dag.Directory()).
		WithExec([]string{"cp", "target/release/" + pkg, resultDir}).
		File(resultDir + "/" + pkg), nil
}
