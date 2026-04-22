package main

import (
	"fmt"

	"dagger/bun/internal/dagger"
)

type Bun struct {
	BunVersion string // +default="1.3"
}

// Returns a cached Bun builder container
func (m *Bun) Builder(source *dagger.Directory) *dagger.Container {
	source = source.WithoutDirectory("target")

	return dag.Container().
		From(fmt.Sprintf("oven/bun:%s-alpine", m.BunVersion)).
		WithExec([]string{"apk", "update"}).
		WithDirectory("/src", source).
		WithWorkdir("/src")
}
