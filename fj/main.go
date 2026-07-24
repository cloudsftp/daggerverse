package main

import (
	"context"
	"fmt"
	"time"

	"dagger/forgejo-release/internal/dagger"
)

type Fj struct {
	// Forgejo CLI Version.
	//
	// +private
	Version string

	// Forgejo token.
	//
	// +private
	Token *dagger.Secret

	// Forgejo host.
	//
	// +private
	Host string
}

func New(
	// Forgejo CLI version.
	//
	// +default="0.6.0"
	version string,

	// Forgejo token.
	token *dagger.Secret,

	// Forgejo host.
	//
	// +optional
	host string,
) *Fj {
	return &Fj{
		Version: version,
		Token:   token,
		Host:    host,
	}
}

// Run a Forgejo CLI command (accepts a list of arguments without "fj").
func (m *Fj) Exec(
	ctx context.Context,

	// Arguments to pass to Forgejo CLI.
	args []string,
) (*dagger.Container, error) {
	args = append([]string{"fj"}, args...)

	c, err := m.Container(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get conainer: %w", err)
	}

	return c.WithExec(args).Sync(ctx)
}

func (m *Fj) host() string {
	if m.Host != "" {
		return m.Host
	}
	return "codeberg.org"
}

func (m *Fj) Container(
	ctx context.Context,
) (*dagger.Container, error) {
	host := m.host()
	c := dag.
		Container().
		From("codeberg.org/forgejo-contrib/forgejo-cli:" + m.Version)

	token, err := m.Token.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get token: %w", err)
	}

	return c.
		WithExec([]string{
			"fj", "--host", host,
			"auth", "add-token",
		}, dagger.ContainerWithExecOpts{
			Stdin: token,
		}).
		WithEnvVariable("CACHE_BUSTER", time.Now().Format(time.RFC3339Nano)).
		Sync(ctx)
}
