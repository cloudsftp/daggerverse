package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"
)

type MergeDirsTests struct{}

func (m *MergeDirsTests) All(ctx context.Context) error {
	var err error

	err = m.TestMergeDirectories(ctx)
	if err != nil {
		return fmt.Errorf("merge directories test failed: %w", err)
	}

	return nil
}

func (m *MergeDirsTests) TestMergeDirectories(ctx context.Context) error {
	var err error

	dir1 := dag.Directory().
		WithNewFile("a", "a").
		WithNewFile("b", "b")

	dir2 := dag.Directory().
		WithNewFile("c", "c").
		WithNewFile("d", "d")

	merged := dag.MergeDirs().
		MergeDirectories([]*dagger.Directory{
			dir1,
			dir2,
		})

	err = assertEntries(ctx, merged, []string{"a", "b", "c", "d"})
	if err != nil {
		return err
	}

	return nil
}

func assertEntries(
	ctx context.Context,
	directory *dagger.Directory,
	expectedFiles []string,
) error {
	entries, err := directory.Entries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get entries: %w", err)
	}

	expectedEntries := map[string]bool{}
	for _, file := range expectedFiles {
		expectedEntries[file] = false
	}

	for _, entry := range entries {
		if _, exists := expectedEntries[entry]; exists {
			expectedEntries[entry] = true
		} else {
			return fmt.Errorf("unexpected file in merged directory: %s", entry)
		}
	}

	for file, found := range expectedEntries {
		if !found {
			return fmt.Errorf("expected file not found in merged directory: %s", file)
		}
	}
	return nil
}
