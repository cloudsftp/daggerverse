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
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := m.TestVet(ctx, source); err != nil {
		return err
	}

	if err := m.TestLint(ctx, source); err != nil {
		return err
	}

	if err := m.TestTest(ctx, source); err != nil {
		return err
	}

	if err := m.TestBuild(ctx, source); err != nil {
		return err
	}

	return nil
}

// Test vetting go code
func (m *GoTests) TestVet(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Go().Vet(ctx, source); err == nil {
		return fmt.Errorf("expected . to fail vet, but it succeeded")
	}

	if err := dag.Go().Vet(ctx, source, dagger.GoVetOpts{Path: "a"}); err != nil {
		return fmt.Errorf("expected a to pass vet: %w", err)
	}

	if err := dag.Go().Vet(ctx, source, dagger.GoVetOpts{Path: "b"}); err == nil {
		return fmt.Errorf("expected b to fail vet, but it succeeded")
	}

	return nil
}

// Test linting go code
func (m *GoTests) TestLint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Go().Lint(ctx, source); err == nil {
		return fmt.Errorf("expected . to fail lint, but it succeeded")
	}

	if err := dag.Go().Lint(ctx, source, dagger.GoLintOpts{Path: "a"}); err != nil {
		return fmt.Errorf("expected a to pass lint: %w", err)
	}

	if err := dag.Go().Lint(ctx, source, dagger.GoLintOpts{Path: "b"}); err == nil {
		return fmt.Errorf("expected b to fail lint, but it succeeded")
	}

	return nil
}

// Test running tests
func (m *GoTests) TestTest(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Go().Test(ctx, source); err == nil {
		return fmt.Errorf("expected . to fail tests, but it succeeded")
	}

	if err := dag.Go().Test(ctx, source, dagger.GoTestOpts{Path: "a"}); err != nil {
		return fmt.Errorf("expected a to pass tests: %w", err)
	}

	if err := dag.Go().Test(ctx, source, dagger.GoTestOpts{Path: "b"}); err == nil {
		return fmt.Errorf("expected b to fail tests, but it succeeded")
	}

	return nil
}

// Test building go executables
func (m *GoTests) TestBuild(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if _, err := dag.Go().Compile(source, "root").Sync(ctx); err != nil {
		return fmt.Errorf("expected root to build: %w", err)
	}

	if _, err := dag.Go().Compile(source, "a", dagger.GoCompileOpts{Path: "a"}).Sync(ctx); err != nil {
		return fmt.Errorf("expected a to build: %w", err)
	}

	if _, err := dag.Go().Compile(source, "b", dagger.GoCompileOpts{Path: "b"}).Sync(ctx); err == nil {
		return fmt.Errorf("expected b to fail build, but it succeeded")
	}

	return nil
}
