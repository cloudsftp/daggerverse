package main

import (
	"context"

	"dagger/renovate/internal/dagger"
)

func (m *Renovate) Container(ctx context.Context) *dagger.Container {
	packages := []string{"npm", "git"}
	packages = append(packages, m.Languages...)

	c := dag.Alpine(dagger.AlpineOpts{
		AlpineVersion: m.AlpineVersion,
		Packages:      packages,
	}).Container().WithExec([]string{
		"npm", "install", "-g", "renovate",
	})

	return c
}
