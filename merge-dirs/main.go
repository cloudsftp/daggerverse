package main

import (
	"context"
	"fmt"
	pathlib "path"

	"dagger/merge-dirs/internal/dagger"
)

type MergeConflictStrategy string

const (
	ErrorOnConflict MergeConflictStrategy = "ERROR"
	KeepRight       MergeConflictStrategy = "KEEP_RIGHT"
	KeepLeft        MergeConflictStrategy = "KEEP_LEFT"
)

type MergeDirs struct{}

// Merge merges multiple directories into one
func (m *MergeDirs) Merge(
	ctx context.Context,
	dirs []*dagger.Directory,
	// +default="ERROR"
	// Conflict resolution strategy: KEEP_LEFT, KEEP_RIGHT, or ERROR (default)
	strategy MergeConflictStrategy,
) (*dagger.Directory, error) {
	if len(dirs) < 2 {
		return nil, fmt.Errorf("need at least 2 directories to merge, got %d", len(dirs))
	}

	left, rest := dirs[0], dirs[1:]

	for i, right := range rest {
		var err error
		left, err = mergeTwoDirectories(ctx, left, right, strategy, "")
		if err != nil {
			return nil, fmt.Errorf("could not merge directory %d: %w", i+1, err)
		}
	}

	return left, nil
}

// mergeTwoDirectories recursively merges entries from right into left
func mergeTwoDirectories(
	ctx context.Context,
	left *dagger.Directory,
	right *dagger.Directory,
	strategy MergeConflictStrategy,
	currentPath string,
) (*dagger.Directory, error) {
	entries, err := right.Entries(ctx, dagger.DirectoryEntriesOpts{
		Path: currentPath,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"could not get entries from right directory at path '%s': %w",
			currentPath, err,
		)
	}

	for _, entry := range entries {
		path := pathlib.Join(currentPath, entry)

		left, err = mergeEntry(ctx, left, right, strategy, path)
		if err != nil {
			return nil, fmt.Errorf("could not merge entry at path '%s': %w", path, err)
		}
	}

	return left, nil
}

// mergeEntry merges a single entry from right into left
func mergeEntry(
	ctx context.Context,
	left *dagger.Directory,
	right *dagger.Directory,
	strategy MergeConflictStrategy,
	path string,
) (*dagger.Directory, error) {
	fileTypeRight, err := right.Stat(path).FileType(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get file type of path '%s': %w", path, err)
	}

	if err := assertAllowedFileType(fileTypeRight); err != nil {
		return nil, fmt.Errorf("right path '%s' has unsupported file type: %w", path, err)
	}

	existsLeft, err := left.Exists(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("could not check whether '%s' exists: %w", path, err)
	}

	if existsLeft {
		return resolveConflict(ctx, left, right, strategy, path, fileTypeRight)
	}

	return copyFromRight(left, right, path, fileTypeRight)
}

// resolveConflict handles merging when both sides have an entry
func resolveConflict(
	ctx context.Context,
	left *dagger.Directory,
	right *dagger.Directory,
	strategy MergeConflictStrategy,
	path string,
	fileTypeRight dagger.FileType,
) (*dagger.Directory, error) {
	fileTypeLeft, err := left.Stat(path).FileType(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get file type of path '%s': %w", path, err)
	}

	if err := assertAllowedFileType(fileTypeLeft); err != nil {
		return nil, fmt.Errorf("left path '%s' has unsupported file type: %w", path, err)
	}

	if fileTypeLeft == dagger.FileTypeDirectory && fileTypeRight == dagger.FileTypeDirectory {
		return mergeTwoDirectories(ctx, left, right, strategy, path)
	}

	switch strategy {
	case ErrorOnConflict:
		return nil, fmt.Errorf(
			"conflict for path '%s': left is %s and right is %s",
			path, fileTypeLeft, fileTypeRight,
		)

	case KeepRight:
		return copyFromRight(left, right, path, fileTypeRight)

	case KeepLeft:
		return left, nil

	default:
		return nil, fmt.Errorf("unexpected strategy: %s", strategy)
	}
}

// copyFromRight copies an entry from right to left
func copyFromRight(
	left *dagger.Directory,
	right *dagger.Directory,
	path string,
	fileType dagger.FileType,
) (*dagger.Directory, error) {
	switch fileType {
	case dagger.FileTypeDirectory:
		directory := right.Directory(path)
		return left.WithDirectory(path, directory), nil

	case dagger.FileTypeRegular:
		file := right.File(path)
		return left.WithFile(path, file), nil

	default:
		return nil, fmt.Errorf(
			"unexpected file type of path '%s': %s",
			path, fileType,
		)
	}
}

func assertAllowedFileType(fileType dagger.FileType) error {
	switch fileType {
	case dagger.FileTypeDirectory, dagger.FileTypeRegular:
		return nil
	case dagger.FileTypeSymlink:
		return fmt.Errorf("symlink not supported")
	case dagger.FileTypeUnknown:
		return fmt.Errorf("unknown file type")
	default:
		return fmt.Errorf("unexpected file type: %s", fileType)
	}
}
