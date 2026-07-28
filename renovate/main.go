package main

import (
	"context"
	"fmt"

	"dagger/renovate/internal/dagger"
)

type Renovate struct {
	Version      string
	Platform     string
	Endpoint     string
	Repositories []string
	ForgeToken   *dagger.Secret
	GithubToken  *dagger.Secret
}

func New(
	// +default="43.278"
	version string,
	// +default="forgejo"
	platform string,
	// +default="https://codeberg.org/api/v1/"
	endpoint string,
	repositories []string,
	forgeToken,
	githubToken *dagger.Secret,
) *Renovate {
	return &Renovate{
		version,
		platform,
		endpoint,
		repositories,
		forgeToken,
		githubToken,
	}
}

func (m *Renovate) Run(
	ctx context.Context,
) (string, error) {
	config := `
		module.exports = {
		  "repositories": [`

	for _, repo := range m.Repositories {
		config += fmt.Sprintf(`"%s", `, repo)
	}

	config += `]
		}`

	return dag.Container().
		From("renovate/renovate:"+m.Version).
		WithNewFile("config.js", config).
		WithSecretVariable("RENOVATE_TOKEN", m.ForgeToken).
		WithSecretVariable("GITHUB_COM_TOKEN", m.GithubToken).
		WithExec([]string{
			"renovate",
			"--platform", m.Platform,
			"--endpoint", m.Endpoint,
		}).
		Stdout(ctx)
}
