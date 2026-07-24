package main

import (
	"context"
	"fmt"
	"time"

	"dagger/forgejo-release/internal/dagger"
)

const (
	fjImage  = "codeberg.org/forgejo-contrib/forgejo-cli"
	fjDigest = "sha256:309b8759b7107a0da2b76b9b9be77c87464f9be9d4d1b6b4d1f9e22ba1145598" // v0.6.0
)

type Fj struct {
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
	// Forgejo token.
	token *dagger.Secret,

	// Forgejo host.
	//
	// +optional
	host string,
) *Fj {
	return &Fj{
		Token: token,
		Host:  host,
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
		From(fjImage + "@" + fjDigest)

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
