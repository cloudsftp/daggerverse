package main

import (
	"dagger/rust/internal/dagger"
)

const resultDir = "/result"

// Build a service executable
func (m *Rust) BuildExecutable(
	source *dagger.Directory,
	pkg string,
) *dagger.File {
	buildCommand := []string{
		"cargo", "build", "--release",
		"-p", pkg,
	}

	return m.Builder(source).
		WithExec(buildCommand).
		WithDirectory(resultDir, dag.Directory()).
		WithExec([]string{"cp", "target/release/" + pkg, resultDir}).
		File(resultDir + "/" + pkg)
}
