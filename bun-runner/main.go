package main

import (
	"fmt"

	"dagger/bun-runner/internal/dagger"
)

type BunRunner struct {
	BunVersion string // +default="1.3"
	AlpineVersion string // +default="3.23"
}

// Returns a cached Bun builder container
func (m *BunRunner) Builder(source *dagger.Directory) *dagger.Container {
	source = source.WithoutDirectory("target")

	return dag.Container().
		From(fmt.Sprintf("oven/bun:%s-alpine", m.BunVersion)).
		WithExec([]string{"apk", "update"}).
		WithDirectory("/src", source).
		WithWorkdir("/src")
}