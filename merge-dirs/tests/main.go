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
	dir1 := dag.Directory().
		WithNewFile("a", "a").
		WithNewFile("b", "b")

	dir2 := dag.Directory().
		WithNewFile("c", "c").
		WithNewFile("d", "d")

	merged := dag.MergeDirs().
		Merge([]*dagger.Directory{dir1, dir2})

	return assertEntries(
		ctx,
		merged,
		[]ExpectedFile{
			{
				name:    "a",
				content: "a",
			},
			{
				name:    "b",
				content: "b",
			},
			{
				name:    "c",
				content: "c",
			},
			{
				name:    "d",
				content: "d",
			},
		},
	)
}

type ExpectedFile struct {
	name    string
	content string
}

func assertEntries(
	ctx context.Context,
	directory *dagger.Directory,
	expectedFiles []ExpectedFile,
) error {
	entries, err := directory.Entries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get entries: %w", err)
	}

	expectedEntries := map[string]bool{}
	expectedContent := map[string]string{}
	for _, file := range expectedFiles {
		expectedEntries[file.name] = false
		expectedContent[file.name] = file.content
	}

	for _, entry := range entries {
		_, ok := expectedEntries[entry]
		if !ok {
			return fmt.Errorf("unexpected file in merged directory: %s", entry)
		}

		content, err := directory.File(entry).Contents(ctx)
		if err != nil {
			return fmt.Errorf("could not get content of merged file '%s': %w", entry, err)
		}

		if content != expectedContent[entry] {
			return fmt.Errorf(
				"file '%s' did not match the expected content: '%s' (expected) != '%s' (actual)",
				entry, expectedContent[entry], content,
			)
		}

		expectedEntries[entry] = true
	}

	for file, found := range expectedEntries {
		if !found {
			return fmt.Errorf("expected file not found in merged directory: %s", file)
		}
	}
	return nil
}
