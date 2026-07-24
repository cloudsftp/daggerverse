package main

import (
	"context"
	"fmt"
	"path/filepath"

	"dagger/forgejo-release/internal/dagger"
)

// Manage releases.
func (m *Fj) Release() *Release {
	return &Release{Fj: m}
}

type Release struct {
	// +private
	Fj *Fj
}

// Create a new GitHub Release for a repository.
func (m *Release) Create(
	ctx context.Context,

	// GitHub repository (e.g. "owner/repo").
	repo string,

	// Release title.
	title string,

	// Tag this release should point to or create.
	tag string,

	// Release assets to upload.
	//
	// +optional
	files []*dagger.File,

	// Save the release as a draft instead of publishing it.
	//
	// +optional
	draft bool,

	// Mark the release as a prerelease.
	//
	// +optional
	preRelease bool,

	// Release notes.
	//
	// +optional
	body string,
) error {
	c, err := m.Fj.Container(ctx)
	if err != nil {
		return fmt.Errorf("could not get fj container: %w", err)
	}

	host := m.Fj.host()

	args := []string{
		"fj",
		"--host", host,
		"release", "create",
		title,
		"--repo", repo,
		"--tag", tag,
	}

	if body != "" {
		args = append(args, "--body", body)
	}

	if draft {
		args = append(args, "--draft")
	}

	if preRelease {
		args = append(args, "--prerelease")
	}

	{
		dir := dag.Directory().WithFiles("", files)

		entries, err := dir.Entries(ctx)
		if err != nil {
			return fmt.Errorf("could not get artifact directory entries: %w", err)
		}

		c = c.WithMountedDirectory("/work/assets", dir)

		for _, e := range entries {
			args = append(
				args,
				"--attach",
				filepath.Join("/work/assets", e),
			)
		}
	}

	_, err = c.WithExec(args).Sync(ctx)

	return err
}
