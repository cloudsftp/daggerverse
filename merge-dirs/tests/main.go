package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/tests/internal/dagger"
)

type MergeDirsTests struct{}

type ExpectedFile struct {
	path    string
	content string
}

type MergeDirectoriesTestCase struct {
	name     string
	dirs     []*dagger.Directory
	strategy dagger.MergeDirsMergeConflictStrategy
	expected []*ExpectedFile
}

type MergeDirectoriesErrorTestCase struct {
	name           string
	dirs           []*dagger.Directory
	strategy       dagger.MergeDirsMergeConflictStrategy
	expectedErrMsg string
}

func (m *MergeDirsTests) All(ctx context.Context) error {
	if err := m.runSuccessTests(ctx); err != nil {
		return fmt.Errorf("success test cases failed: %w", err)
	}

	if err := m.runErrorTests(ctx); err != nil {
		return fmt.Errorf("success test cases failed: %w", err)
	}

	return nil
}

func (m *MergeDirsTests) runSuccessTests(ctx context.Context) error {
	tests := []MergeDirectoriesTestCase{
		// Basic merge tests
		{
			name: "disjunct files at root",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"a": "a", "b": "b"}),
				buildDirectory(map[string]any{"c": "c", "d": "d"}),
			},
			expected: []*ExpectedFile{
				{path: "a", content: "a"},
				{path: "b", content: "b"},
				{path: "c", content: "c"},
				{path: "d", content: "d"},
			},
		},

		// File conflict resolution tests
		{
			name: "file conflict: keep left",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"a": "a", "b": "b1"}),
				buildDirectory(map[string]any{"b": "b2", "c": "c"}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			expected: []*ExpectedFile{
				{path: "a", content: "a"},
				{path: "b", content: "b1"},
				{path: "c", content: "c"},
			},
		},
		{
			name: "file conflict: keep right",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"a": "a", "b": "b1"}),
				buildDirectory(map[string]any{"b": "b2", "c": "c"}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			expected: []*ExpectedFile{
				{path: "a", content: "a"},
				{path: "b", content: "b2"},
				{path: "c", content: "c"},
			},
		},

		// Nested directory tests
		{
			name: "merge nested directories",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"dir": map[string]any{"a": "a"}}),
				buildDirectory(map[string]any{"dir": map[string]any{"b": "b"}}),
			},
			expected: []*ExpectedFile{
				{path: "dir/a", content: "a"},
				{path: "dir/b", content: "b"},
			},
		},
		{
			name: "deep nesting (3 levels)",
			dirs: []*dagger.Directory{
				dag.Directory().WithDirectory("level1",
					dag.Directory().WithDirectory("level2",
						dag.Directory().WithNewFile("deep", "deep1"))),
				dag.Directory().WithDirectory("level1",
					dag.Directory().WithDirectory("level2",
						dag.Directory().WithNewFile("other", "other2"))),
			},
			expected: []*ExpectedFile{
				{path: "level1/level2/deep", content: "deep1"},
				{path: "level1/level2/other", content: "other2"},
			},
		},

		// Multiple directory merge
		{
			name: "merge three directories with chained conflicts",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"a": "a1", "b": "b1"}),
				buildDirectory(map[string]any{"b": "b2", "c": "c2"}),
				buildDirectory(map[string]any{"c": "c3", "d": "d3"}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			expected: []*ExpectedFile{
				{path: "a", content: "a1"},
				{path: "b", content: "b2"},
				{path: "c", content: "c3"},
				{path: "d", content: "d3"},
			},
		},

		// Nested file conflicts
		{
			name: "nested file conflict: keep left",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"dir": map[string]any{"conflict": "left"}}),
				buildDirectory(map[string]any{"dir": map[string]any{"conflict": "right"}}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			expected: []*ExpectedFile{
				{path: "dir/conflict", content: "left"},
			},
		},
		{
			name: "nested file conflict: keep right",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"dir": map[string]any{"conflict": "left"}}),
				buildDirectory(map[string]any{"dir": map[string]any{"conflict": "right"}}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			expected: []*ExpectedFile{
				{path: "dir/conflict", content: "right"},
			},
		},

		// Type mismatch tests: directory vs file
		{
			name: "type mismatch (dir left, file right): keep left",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"mismatch": map[string]any{"inner": "inner-content"}}),
				buildDirectory(map[string]any{"mismatch": "file-content"}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			expected: []*ExpectedFile{
				{path: "mismatch/inner", content: "inner-content"},
			},
		},
		{
			name: "type mismatch (dir left, file right): keep right",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"mismatch": map[string]any{"inner": "inner-content"}}),
				buildDirectory(map[string]any{"mismatch": "file-content"}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			expected: []*ExpectedFile{
				{path: "mismatch", content: "file-content"},
			},
		},
		{
			name: "type mismatch (file left, dir right): keep left",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"mismatch": "file-content"}),
				buildDirectory(map[string]any{"mismatch": map[string]any{"inner": "inner-content"}}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepLeft,
			expected: []*ExpectedFile{
				{path: "mismatch", content: "file-content"},
			},
		},
		{
			name: "type mismatch (file left, dir right): keep right",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"mismatch": "file-content"}),
				buildDirectory(map[string]any{"mismatch": map[string]any{"inner": "inner-content"}}),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			expected: []*ExpectedFile{
				{path: "mismatch/inner", content: "inner-content"},
			},
		},
	}

	for _, test := range tests {
		err := test.run(ctx)
		if err != nil {
			return fmt.Errorf("test case '%s' failed: %w", test.name, err)
		}
	}

	return nil
}

func (m *MergeDirsTests) runErrorTests(ctx context.Context) error {
	tests := []MergeDirectoriesErrorTestCase{
		{
			name: "type mismatch (dir left, file right): error strategy",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"mismatch": map[string]any{"inner": "inner-content"}}),
				buildDirectory(map[string]any{"mismatch": "file-content"}),
			},
			expectedErrMsg: "conflict",
		},
		{
			name: "type mismatch (file left, dir right): error strategy",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"mismatch": "file-content"}),
				buildDirectory(map[string]any{"mismatch": map[string]any{"inner": "inner-content"}}),
			},
			expectedErrMsg: "conflict",
		},
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

func (t *MergeDirectoriesErrorTestCase) run(ctx context.Context) error {
	merged := dag.MergeDirs().Merge(
		t.dirs,
		dagger.MergeDirsMergeOpts{
			Strategy: t.strategy,
		},
	)

	_, err := merged.Entries(ctx, dagger.DirectoryEntriesOpts{})
	if err == nil {
		return fmt.Errorf(
			"expected error containing '%s' but got none",
			t.expectedErrMsg,
		)
	}

	if !strings.Contains(err.Error(), t.expectedErrMsg) {
		return fmt.Errorf(
			"expected error containing '%s' but got: %v",
			t.expectedErrMsg, err,
		)
	}

	return nil
}

func assertEntries(
	ctx context.Context,
	directory *dagger.Directory,
	expectedFiles []*ExpectedFile,
) error {
	for _, expectedFile := range expectedFiles {
		exists, err := directory.Exists(ctx, expectedFile.path, dagger.DirectoryExistsOpts{
			ExpectedType: dagger.ExistsTypeRegularType,
		})
		if err != nil {
			return fmt.Errorf(
				"could not check, whether '%s' exists: %w",
				expectedFile.path, err,
			)
		}
		if !exists {
			return fmt.Errorf("file '%s' does not exist", expectedFile.path)
		}

		content, err := directory.File(expectedFile.path).Contents(ctx)
		if err != nil {
			return fmt.Errorf(
				"could not read file at path '%s': %w",
				expectedFile.path, err,
			)
		}

		if content != expectedFile.content {
			return fmt.Errorf(
				"unexpected content in file '%s': expected '%s' != got '%s'",
				expectedFile.path, expectedFile.content, content,
			)
		}
	}

	allFilePaths, err := collectAllFiles(ctx, directory, "", []string{})
	if err != nil {
		return fmt.Errorf("could not get all file paths: %w", err)
	}

	expectedFilePathsSet := make(map[string]struct{}, len(expectedFiles))
	for _, expectedFile := range expectedFiles {
		expectedFilePathsSet[expectedFile.path] = struct{}{}
	}

	for _, filePath := range allFilePaths {
		if _, ok := expectedFilePathsSet[filePath]; !ok {
			return fmt.Errorf("unexpected file path: %s", filePath)
		}
	}

	return nil
}
