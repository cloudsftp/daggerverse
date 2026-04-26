package main

import (
	"context"

	"dagger/tests/internal/dagger"
)

type RustTests struct{}

// Run all rust tests
func (m *RustTests) Run(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := m.TestBuild(ctx, source); err != nil {
		return err
	}

	return nil
}

// Test checkint rust code
func (m *RustTests) TestCheck(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Rust().Check(ctx, source); err != nil {
		return err
	}

	return nil
}

// Test run tests
func (m *RustTests) TestTest(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	if err := dag.Rust().Test(ctx, source); err != nil {
		return err
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
		return err
	}

	if _, err := dag.Rust().BuildExecutable(source, "bin-a").Sync(ctx); err != nil {
		return err
	}

	if _, err := dag.Rust().BuildExecutable(source, "bin-b").Sync(ctx); err != nil {
		return err
	}

	return nil
}
