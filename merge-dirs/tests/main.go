package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"
)

type MergeDirsTests struct{}

type ExpectedFile struct {
	name    string
	content string
}

type MergeDirectoriesTestCase struct {
	name     string
	dirs     []*dagger.Directory
	strategy dagger.MergeDirsMergeConflictStrategy
	expected []*ExpectedFile
}

func (m *MergeDirsTests) All(ctx context.Context) error {
	tests := []MergeDirectoriesTestCase{
		{
			name: "disjunct files root",
			dirs: []*dagger.Directory{
				dag.Directory().
					WithNewFile("a", "a").
					WithNewFile("b", "b"),
				dag.Directory().
					WithNewFile("c", "c").
					WithNewFile("d", "d"),
			},
			expected: []*ExpectedFile{
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
		},
		{
			name: "conflicting files root, left",
			dirs: []*dagger.Directory{
				dag.Directory().
					WithNewFile("a", "a").
					WithNewFile("b", "b1"),
				dag.Directory().
					WithNewFile("b", "b2").
					WithNewFile("c", "c"),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			expected: []*ExpectedFile{
				{
					name:    "a",
					content: "a",
				},
				{
					name:    "b",
					content: "b1",
				},
				{
					name:    "c",
					content: "c",
				},
			},
		},
		{
			name: "conflicting files root, right",
			dirs: []*dagger.Directory{
				dag.Directory().
					WithNewFile("a", "a").
					WithNewFile("b", "b1"),
				dag.Directory().
					WithNewFile("b", "b2").
					WithNewFile("c", "c"),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			expected: []*ExpectedFile{
				{
					name:    "a",
					content: "a",
				},
				{
					name:    "b",
					content: "b2",
				},
				{
					name:    "c",
					content: "c",
				},
			},
		},
		/*
			{
				name: "merge nested directory",
				dirs: []*dagger.Directory{
					dag.Directory().
						WithNewDirectory("dir").
						WithNewFile("a", "a"),
					dag.Directory().
						WithNewDirectory("dir").
						WithNewFile("b", "b"),
				},
				strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
				expected: []*ExpectedFile{
					{
						name:    "dir/a",
						content: "a",
					},
					{
						name:    "dir/b",
						content: "b",
					},
				},
			},
		*/
	}

	for _, test := range tests {
		err := test.run(ctx)
		if err != nil {
			return fmt.Errorf("test case '%s' failed: %w", test.name, err)
		}
	}

	return nil
}

func (t *MergeDirectoriesTestCase) run(ctx context.Context) error {
	merged := dag.MergeDirs().Merge(
		t.dirs,
		dagger.MergeDirsMergeOpts{
			Strategy: t.strategy,
		},
	)
	return assertEntries(ctx, merged, t.expected)
}

func assertEntries(
	ctx context.Context,
	directory *dagger.Directory,
	expectedFiles []*ExpectedFile,
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
