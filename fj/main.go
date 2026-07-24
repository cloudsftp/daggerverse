package main

import (
	"context"
	"fmt"

	"dagger/forgejo-release/internal/dagger"
)

type Fj struct {
	// Forgejo token.
	//
	// +private
	Token *dagger.Secret

	// Forgejo repository (e.g. "owner/repo").
	//
	// +private
	Repository string

	// Forgejo host.
	//
	// +private
	Host string

	// Forgejo repository source (with .git directory).
	Source *dagger.Directory
}

func New(
	// GitHub token.
	//
	// +optional
	token *dagger.Secret,

	// GitHub repository (e.g. "owner/repo").
	//
	// +optional
	repo string,

	// GitHub host.
	//
	// +optional
	host string,

	// Git repository source (with .git directory).
	//
	// +optional
	source *dagger.Directory,
) *Fj {
	return &Fj{
		Token:      token,
		Repository: repo,
		Host:       host,
		// CACert:     caCert,
		Source: source,
	}
}

// Run a Forgejo CLI command (accepts a list of arguments without "fj").
func (m *Fj) Exec(
	ctx context.Context,

	// Arguments to pass to Forgejo CLI.
	args []string,
) (*dagger.Container, error) {
	args = append([]string{"gh"}, args...)

	ctr, err := m.container(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get conainer: %w", err)
	}

	return ctr.WithExec(args).Sync(ctx)
}

const (
	fjImage = "codeberg.org/forgejo-contrib/forgejo-cli:latest"
)

func (m *Fj) host() string {
	if m.Host != "" {
		return m.Host
	}
	return "codeberg.org"
}

func (m *Fj) container(
	ctx context.Context,
) (*dagger.Container, error) {
	host := m.host()
	ctr := dag.
		Container().
		From(fjImage)

	token, err := m.Token.Plaintext(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get token: %w", err)
	}

	return ctr.
		WithExec([]string{
			"fj", "--host", host,
			"auth", "add-token",
		}, dagger.ContainerWithExecOpts{
			Stdin: token,
		}).
		With(func(c *dagger.Container) *dagger.Container {
			if m.Source != nil {
				c = c.
					WithWorkdir("/work/repo").
					WithMountedDirectory("/work/repo", m.Source)
			}

			return c
		}).Sync(ctx)
}

func (m *Fj) Fj(
	ctx context.Context,
) (*dagger.Container, error) {
	return m.container(ctx)
}

// Create forgejo release
func (m *Fj) CreateRelease(
	ctx context.Context,
	tag string,
	// +optional
	assets *dagger.Directory,
) error {
	c, err := m.Fj(ctx)
	if err != nil {
		return err
	}

	_, err = c.
		WithExec([]string{"fj", "release"}).
		Sync(ctx)

	return err
}

/*
func (m *ForgejoRelease) CreateRelease(
	ctx context.Context,
	tag string,
	// +optional
	assets *dagger.Directory,
) error {
	title := tag

	err := dag.Gh(dagger.GhOpts{
		Token: m.Token,
		Repo:  m.Repository,
		Host:  m.Forge,
	}).
		Release().
		Create(ctx, tag, title)

	return err
}
*/

/*
func (m *ForgejoRelease) CreateRelease(
	ctx context.Context,
	tag string,
	// +optional
	assets *dagger.Directory,
) error {
	c := dag.Alpine(dagger.AlpineOpts{
		AlpineVersion: m.AlpineVersion,
		Packages:      []string{"curl"},
	}).Container()

	apiEndpoint := m.Forge + "/api/v1/"
	method := "POST"
	path := "repos/" + m.Repository + "/releases"

	tokenPlain, err := m.Token.Plaintext(ctx)
	if err != nil {
		return err
	}

	payloadContent := map[string]string{
		"tag_name": tag,
	}
	payload, err := json.Marshal(payloadContent)
	if err != nil {
		return fmt.Errorf("could not encode payload content: %w", err)
	}

	_, err = c.
		WithExec([]string{
			"curl", "--retry", "5", "--fail",
			"-X", method, "-sS", "-H", "Authorization: token " + tokenPlain,
			"-H", "Content-Type: application/json",
			"-d", string(payload),
			apiEndpoint + path,
		}).
		Sync(ctx)

	return err
}
*/
