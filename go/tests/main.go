package main

import (
	"context"
	"fmt"

	"dagger/go-tests/internal/dagger"
)

type GoTests struct{}

// Run all go tests
func (m *GoTests) Run(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	if err := m.TestLint(ctx, source); err != nil {
		return err
	}

	if err := m.TestBuildExecutable(ctx, source); err != nil {
		return err
	}

	if err := m.TestBuildImage(ctx, source); err != nil {
		return err
	}

	return nil
}

// Test linting go code
func (m *GoTests) TestLint(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	err := dag.Go().Lint(ctx, source)
	if err != nil {
		return fmt.Errorf("could not lint: %w", err)
	}

	return nil
}

// Test building go executables
func (m *GoTests) TestBuildExecutable(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	_, err := dag.Go().Compile(source, "executable").Sync(ctx)
	if err != nil {
		return fmt.Errorf("could not build executable: %w", err)
	}

	return nil
}

// Test building the server image
func (m *GoTests) TestBuildImage(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	_, err := dag.Go().BuildImage(source, "executable").Sync(ctx)
	if err != nil {
		return fmt.Errorf("could not build image: %w", err)
	}

	return nil
}
