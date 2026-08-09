package main

import (
	"dagger/alpine/internal/dagger"
)

const (
	user = "appuser"
)

type Alpine struct {
	AlpineVersion string
	Packages      []string
}

func New(
	// +default="3.24"
	alpineVersion string,
	// +optional
	packages []string,
) *Alpine {
	return &Alpine{
		alpineVersion,
		packages,
	}
}

// Returns a bare alpine container
func (m *Alpine) Container() *dagger.Container {
	container := dag.Container().From("alpine:" + m.AlpineVersion)

	if len(m.Packages) > 0 {
		container = container.WithExec(append(
			[]string{"apk", "add", "--no-cache"},
			m.Packages...,
		))
	}

	return container
}

// Create a minimal service container from an executable
func (m *Alpine) ServiceContainer(
	executable *dagger.File,
	name string,
	// +optional
	platform string,
) *dagger.Container {
	executablePath := "/" + name

	opts := dagger.ContainerOpts{}
	if platform != "" {
		opts.Platform = dagger.Platform(platform)
	}

	return dag.Container(opts).
		From("alpine:"+m.AlpineVersion).
		WithLabel("", "").
		WithExec([]string{"adduser", user, "-D"}).
		WithUser(user).
		WithFile(executablePath, executable).
		WithEntrypoint([]string{executablePath})
}
