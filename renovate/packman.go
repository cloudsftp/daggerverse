package main

import (
	"context"
	"fmt"

	"dagger/renovate/internal/dagger"
)

type installFunc func(*dagger.Container) *dagger.Container

func (m *Renovate) Container(ctx context.Context) (*dagger.Container, error) {
	packages := []string{"npm", "git"}
	installationSteps := []installFunc{}

	for _, lang := range m.Languages {
		switch lang {
		case "cargo":
			packages = append(packages, "cargo")
		case "go":
			packages = append(packages, "go")
		case "uv":
			packages = append(packages, "uv")
		case "bun":
			installationSteps = append(installationSteps,
				func(c *dagger.Container) *dagger.Container {
					return c.WithExec([]string{"npm", "install", "-g", "bun"})
				},
			)
		default:
			return nil, fmt.Errorf("unsupported language/tool: %s", lang)
		}
	}

	c := dag.Alpine(dagger.AlpineOpts{
		AlpineVersion: m.AlpineVersion,
		Packages:      packages,
	}).Container().WithExec([]string{
		"npm", "install", "-g", "renovate",
	})

	for _, step := range installationSteps {
		c = step(c)
	}

	return c, nil
}
