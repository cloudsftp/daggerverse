package main

import (
	"strings"

	"dagger/go/internal/dagger"
)

// Compile executable
func (m *Go) Compile(
	source *dagger.Directory,
	name string,
	// +default=""
	path string,
	// +default="linux"
	os string,
	// +default="amd64"
	arch string,
	// +optional
	tags []string,
	// +optional
	ldflags string,
) *dagger.File {
	executablePath := "/tmp/" + name

	args := []string{
		"go", "build",
	}

	if len(tags) > 0 {
		args = append(
			args,
			"-tags",
			strings.Join(tags, " "),
		)
	}

	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}

	args = append(
		args,
		"-o", executablePath,
		resolvePath(path),
	)

	return m.Builder(source).
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithExec(args).
		File(executablePath)
}
