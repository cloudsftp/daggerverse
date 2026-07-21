package main

import (
	"context"
	"dagger/git-pages/internal/dagger"
	"fmt"
)

type GitPages struct {
	GitPagesVersion string
}

func New(
	// +default="1.10.0"
	gitPagesVersion string,
) *GitPages {
	return &GitPages{
		gitPagesVersion,
	}
}

// Deploys a git page
func (m *GitPages) Deploy(
	ctx context.Context,
	dist *dagger.Directory,
	token *dagger.Secret,
	site, server string,
) error {
	g := dag.Container().
		From("codeberg.org/git-pages/git-pages-cli:" + m.GitPagesVersion)

	tokenPlain, err := token.Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("could not get token contents: %w", err)
	}

	_, err = g.
		WithDirectory("/dist", dist).
		WithSecretVariable("TOKEN", token).
		WithExec([]string{
			"git-pages-cli",
			site,
			"--token", tokenPlain,
			"--server", server,
			"--upload-dir", "/dist",
		}).
		Sync(ctx)

	return err
}
