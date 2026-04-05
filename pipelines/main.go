package main

import (
	"context"
	"fmt"
	pathlib "path"

	"dagger/clouds-dagger-modules/internal/dagger"
)

const (
	golangLintVersion = "2.11.4"
)

type CloudsDaggerModules struct{}

// Run the whole pipeline
func (m CloudsDaggerModules) Run(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) error {
	if err := m.Lint(ctx, source); err != nil {
		return err
	}

	if err := m.Test(ctx); err != nil {
		return err
	}

	return nil
}

// Lint Code
func (m CloudsDaggerModules) Lint(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) error {
	run := func(path string) error {
		source := source.Directory(path)

		_, err := dag.Container().
			From("golangci/golangci-lint:v"+golangLintVersion+"-alpine").
			WithMountedDirectory("/app", source).
			WithWorkdir("/app").
			WithExec([]string{"golangci-lint", "run", "./..."}).
			Sync(ctx)

		if err != nil {
			return fmt.Errorf("lint failed: %w", err)
		}

		return nil
	}

	for _, path := range []string{
		"bun-runner",
		"go-runner",
		"merge-dirs",
		"merge-dirs",
		"pipelines",
		"rust-runner",
	} {
		if err := run(path); err != nil {
			return err
		}

		testsPath := pathlib.Join(path, "tests")
		testsExist, err := source.Exists(ctx, testsPath)
		if err != nil {
			return fmt.Errorf(
				"could not check, whether '%s' exists: %w",
				testsPath, err,
			)
		}

		if testsExist {
			if err := run(testsPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run Tests
func (m *CloudsDaggerModules) Test(ctx context.Context) error {
	if err := dag.MergeDirsTests().All(ctx); err != nil {
		return fmt.Errorf("merge directories tests failed: %w", err)
	}

	return nil
}
