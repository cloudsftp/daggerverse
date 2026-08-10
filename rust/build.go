package main

import (
	"dagger/rust/internal/dagger"
)

const resultDir = "/result"

// Build a service executable
func (m *Rust) BuildExecutable(
	source *dagger.Directory,
	pkg string,
	// +optional
	target string,
	// +default=false
	normal bool,
) *dagger.File {
	buildVerb := ""
	if normal {
		buildVerb = "build"
	} else {
		buildVerb = "zigbuild"
	}

	buildCommand := []string{
		"cargo", buildVerb, "--release",
		"-p", pkg,
	}

	targetBinaryPath := "target/"
	if target != "" {
		buildCommand = append(buildCommand, "--target", target)
		targetBinaryPath += target + "/"
	}
	targetBinaryPath += "release/" + pkg

	return m.Builder(source).
		WithExec(buildCommand).
		WithDirectory(resultDir, dag.Directory()).
		WithExec([]string{"cp", targetBinaryPath, resultDir}).
		File(resultDir + "/" + pkg)
}
