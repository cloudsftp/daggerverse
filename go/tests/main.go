package main

import (
	"context"
	"fmt"

	"dagger/go-tests/internal/dagger"
)

type GoTests struct{}

func (m *GoTests) All(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	if err := m.TestBuildExecutable(ctx, source); err != nil {
		return err
	}

	if err := m.TestBuildImage(ctx, source); err != nil {
		return err
	}

	if err := m.TestLint(ctx, source); err != nil {
		return err
	}

	return nil
}

func (m *GoTests) TestBuildExecutable(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	_, err := dag.Go().
		Compile(dagger.GoCompileOpts{Source: source}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("could not build executable: %w", err)
	}

	return nil
}

func (m *GoTests) TestBuildImage(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	_, err := dag.Go().
		BuildImage(dagger.GoBuildImageOpts{Source: source}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("could not build image: %w", err)
	}

	return nil
}

func (m *GoTests) TestLint(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	_, err := dag.Go().
		Lint(dagger.GoLintOpts{Source: source}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("could not lint: %w", err)
	}

	return nil
}
