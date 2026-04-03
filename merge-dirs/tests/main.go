package main

import (
	"context"

	"dagger/tests/internal/dagger"
)

type Tests struct{}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Tests) GrepDir(ctx context.Context) error {
	var err error
	err = err

	_ = dag.MergeDirs().MergeDirectories([]*dagger.Directory{})

	return nil
}
