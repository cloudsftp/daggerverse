package main

import (
	"context"
	"fmt"

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
					path:    "a",
					content: "a",
				},
				{
					path:    "b",
					content: "b",
				},
				{
					path:    "c",
					content: "c",
				},
				{
					path:    "d",
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
					path:    "a",
					content: "a",
				},
				{
					path:    "b",
					content: "b1",
				},
				{
					path:    "c",
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
					path:    "a",
					content: "a",
				},
				{
					path:    "b",
					content: "b2",
				},
				{
					path:    "c",
					content: "c",
				},
			},
		},
		{
			name: "merge nested directory",
			dirs: []*dagger.Directory{
				dag.Directory().
					WithDirectory(
						"dir",
						dag.Directory().
							WithNewFile("a", "a"),
					),
				dag.Directory().
					WithDirectory(
						"dir",
						dag.Directory().
							WithNewFile("b", "b"),
					),
			},
			strategy: dagger.MergeDirsMergeConflictStrategyKeepRight,
			expected: []*ExpectedFile{
				{
					path:    "dir/a",
					content: "a",
				},
				{
					path:    "dir/b",
					content: "b",
				},
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

	expectedFilePathsSet := map[string]struct{}{}
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

func collectAllFiles(
	ctx context.Context,
	directory *dagger.Directory,
	currentPath string,
	files []string,
) ([]string, error) {
	entries, err := directory.Entries(ctx, dagger.DirectoryEntriesOpts{
		Path: currentPath,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"could not get entries of path '%s': %w",
			currentPath, err,
		)
	}

	for _, entry := range entries {
		path := currentPath + entry

		fileType, err := directory.Stat(path).FileType(ctx)
		if err != nil {
			return nil, fmt.Errorf(
				"could not get file type of path '%s': %w",
				path, err,
			)
		}

		switch fileType {
		case dagger.FileTypeDirectory:
			filesInDirectory, err := collectAllFiles(ctx, directory, path, files)
			if err != nil {
				return nil, err
			}
			files = append(files, filesInDirectory...)

		case dagger.FileTypeRegular:
			files = append(files, path)

		case dagger.FileTypeSymlink:
			return nil, fmt.Errorf(
				"path '%s' is a symlink: %w",
				path, err,
			)

		case dagger.FileTypeUnknown:
			return nil, fmt.Errorf(
				"unknown file type of path '%s': %w",
				path, err,
			)

		default:
			return nil, fmt.Errorf(
				"file type of path '%s' unexpected: %w",
				path, err,
			)
		}
	}

	return files, nil
}
