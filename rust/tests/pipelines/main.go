package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"
)

type RustTests struct{}

// Run all rust tests
func (m *RustTests) Run(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := m.TestCheck(ctx, source); err != nil {
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

// Test checking rust code
func (m *RustTests) TestCheck(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Rust().Check(ctx, source); err == nil {
		return fmt.Errorf("expected . to fail check, but it succeeded")
	}

	if err := dag.Rust().Check(ctx, source, dagger.RustCheckOpts{Pkg: "root"}); err != nil {
		return fmt.Errorf("expected root to pass check: %w", err)
	}

	if err := dag.Rust().Check(ctx, source, dagger.RustCheckOpts{Pkg: "bin-a"}); err != nil {
		return fmt.Errorf("expected bin-a to pass check: %w", err)
	}

	if err := dag.Rust().Check(ctx, source, dagger.RustCheckOpts{Pkg: "bin-b"}); err == nil {
		return fmt.Errorf("expected bin-b to fail check, but it succeeded")
	}

	return nil
}

// Test linting rust code
func (m *RustTests) TestLint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Rust().Lint(ctx, source); err == nil {
		return fmt.Errorf("expected . to fail lint, but it succeeded")
	}

	if err := dag.Rust().Lint(ctx, source, dagger.RustLintOpts{Pkg: "root"}); err != nil {
		return fmt.Errorf("expected root to pass lint: %w", err)
	}

	if err := dag.Rust().Lint(ctx, source, dagger.RustLintOpts{Pkg: "bin-a"}); err != nil {
		return fmt.Errorf("expected bin-a to pass lint: %w", err)
	}

	if err := dag.Rust().Lint(ctx, source, dagger.RustLintOpts{Pkg: "bin-b"}); err == nil {
		return fmt.Errorf("expected bin-b to fail lint, but it succeeded")
	}

	return nil
}

// Test running tests
func (m *RustTests) TestTest(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Rust().Test(ctx, source); err == nil {
		return fmt.Errorf("expected . to fail tests, but it succeeded")
	}

	if err := dag.Rust().Test(ctx, source, dagger.RustTestOpts{Pkg: "root"}); err != nil {
		return fmt.Errorf("expected root to pass tests: %w", err)
	}

	if err := dag.Rust().Test(ctx, source, dagger.RustTestOpts{Pkg: "bin-a"}); err != nil {
		return fmt.Errorf("expected bin-a to pass tests: %w", err)
	}

	if err := dag.Rust().Test(ctx, source, dagger.RustTestOpts{Pkg: "bin-b"}); err == nil {
		return fmt.Errorf("expected bin-b to fail tests, but it succeeded")
	}

	return nil
}

// Test building rust executables
func (m *RustTests) TestBuild(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if _, err := dag.Rust().BuildExecutable(source, "root").Sync(ctx); err != nil {
		return fmt.Errorf("expected root to build: %w", err)
	}

	if _, err := dag.Rust().BuildExecutable(source, "bin-a").Sync(ctx); err != nil {
		return fmt.Errorf("expected bin-a to build: %w", err)
	}

	if _, err := dag.Rust().BuildExecutable(source, "bin-b").Sync(ctx); err == nil {
		return fmt.Errorf("expected bin-b to fail build, but it succeeded")
	}

	return nil
}
