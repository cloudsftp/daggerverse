package main

import (
	"context"
	"time"

	"dagger/renovate/internal/dagger"
)

type Renovate struct {
	AlpineVersion string
	Languages     []string
	Platform      string
	Endpoint      string
}

func New(
	// +default=""
	alpineVersion string,
	// +optional
	languages []string,
	// +default="forgejo"
	platform string,
	// +default="https://codeberg.org/api/v1/"
	endpoint string,
) *Renovate {
	return &Renovate{
		alpineVersion,
		languages,
		platform,
		endpoint,
	}
}

func (m *Renovate) Run(
	ctx context.Context,
	// +default=false
	debugLog bool,
	repositories []string,
	renovateToken,
	githubComToken *dagger.Secret,
) (string, error) {
	cmd := append(
		[]string{
			"npx", "renovate",
			"--platform", m.Platform,
			"--endpoint", m.Endpoint,
		},
		repositories...,
	)

	c := m.Container(ctx).
		WithSecretVariable("RENOVATE_TOKEN", renovateToken).
		WithSecretVariable("GITHUB_COM_TOKEN", githubComToken)

	if debugLog {
		c = c.WithEnvVariable("LOG_LEVEL", "debug")
	}

	return c.
		WithEnvVariable("CACHE_BUSTER", time.Now().Format(time.RFC3339Nano)).
		WithExec(cmd).
		Stdout(ctx)
}
