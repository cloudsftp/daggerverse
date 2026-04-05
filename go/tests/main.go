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

	return nil
}

func (m *GoTests) TestBuildExecutable(
	ctx context.Context,
	// +defaultPath="./data"
	source *dagger.Directory,
) error {
	_, err := dag.Go().
		BuildExecutable(source, "main.go").
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("could not build executable: %w", err)
	}

	return nil
}
