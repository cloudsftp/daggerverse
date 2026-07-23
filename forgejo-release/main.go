package main

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger/forgejo-release/internal/dagger"
)

type ForgejoRelease struct {
	AlpineVersion string

	Forge      string
	Repository string
	Token      *dagger.Secret
}

func New(
	// +default="3.24"
	alpineVersion string,
	forge string,
	repository string,
	token *dagger.Secret,
) *ForgejoRelease {
	return &ForgejoRelease{
		AlpineVersion: alpineVersion,
		Forge:         forge,
		Repository:    repository,
		Token:         token,
	}
}

// Release
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
