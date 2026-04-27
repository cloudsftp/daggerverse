package main

import (
	"dagger/alpine/internal/dagger"
)

const (
	user = "appuser"
)

type Alpine struct {
	AlpineVersion string
}

func New(
	// +default="3.23"
	alpineVersion string,
) *Alpine {
	return &Alpine{
		alpineVersion,
	}
}

// Create a minimal service container from an executable
func (m *Alpine) ServiceContainer(
	executable *dagger.File,
	// +default="server"
	name string,
) *dagger.Container {
	return dag.Container().
		From("alpine:"+m.AlpineVersion).
		WithExec([]string{"adduser", user, "-D"}).
		WithUser(user).
		WithFile("/"+name, executable).
		WithEntrypoint([]string{"/" + name})
}
