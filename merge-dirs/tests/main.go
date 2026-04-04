package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/tests/internal/dagger"
)

type MergeDirsTests struct{}

type MergeDirectoriesTestCase struct {
	name     string
	dirs     []*dagger.Directory
	strategy dagger.MergeDirsMergeConflictStrategy
	expected *dagger.Directory
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
		// Conflict strategy: error
		{
			name: "disjunct files at root",
			dirs: []*dagger.Directory{
				buildDirectory(map[string]any{"a": "a", "b": "b"}),
				buildDirectory(map[string]any{"c": "c", "d": "d"}),
			},
			expected: buildDirectory(map[string]any{"a": "a", "b": "b", "c": "c", "d": "d"}),
		},

		// Conflict strategy: left

		// Conflict strategy: right

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

	return assertDirectory(ctx, t.expected, merged)
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
