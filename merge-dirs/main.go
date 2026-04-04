package main

import (
	"context"
	"fmt"

	"dagger/merge-dirs/internal/dagger"
)

type MergeConflictStrategy string

const (
	KeepLeft        MergeConflictStrategy = "KEEP_LEFT"
	KeepRight       MergeConflictStrategy = "KEEP_RIGHT"
	ErrorOnConflict MergeConflictStrategy = "ERROR"
)

type MergeDirs struct{}

// Merge merges multiple directories into one
func (m *MergeDirs) Merge(
	ctx context.Context,
	dirs []*dagger.Directory,
	// +default="ERROR"
	// Conflict resolution strategy: left (default), right, or error
	strategy MergeConflictStrategy,
) (*dagger.Directory, error) {
	if len(dirs) < 2 {
		return nil, fmt.Errorf("need at least 2 directories to merge, got %d", len(dirs))
	}

	first, rest := dirs[0], dirs[1:]
	var err error

	for i, next := range rest {
		first, err = mergeDirectories2(ctx, first, next, strategy, "")
		if err != nil {
			return nil, fmt.Errorf("could not merge directory %d: %w", i+1, err)
		}
	}

	return first, nil
}

// Merge two directories into one
func mergeDirectories2(
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
		return nil, fmt.Errorf("could not get entries from second directory: %w", err)
	}

	for _, entry := range entries {
		path := currentPath + entry

		left, err = copyPath(ctx, left, right, strategy, path)
		if err != nil {
			return nil, fmt.Errorf("could not copy at path '%s': %w", path, err)
		}
	}

	return left, nil
}

// Copy a specific path from one directory to another
func copyPath(
	ctx context.Context,
	left *dagger.Directory,
	right *dagger.Directory,
	strategy MergeConflictStrategy,
	path string,
) (*dagger.Directory, error) {
	existsLeft, err := left.Exists(ctx, path)
	if err != nil {
		return nil, fmt.Errorf(
			"could not check, whether '%s' exists: %w",
			path, err,
		)
	}

	fileTypeRight, err := right.Stat(path).FileType(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"could not get file type of path '%s': %w",
			path, err,
		)
	}

	if existsLeft {
		return copyPathLeftExists(
			ctx, left, right, strategy,
			path, fileTypeRight,
		)
	}

	err = assertAllowedFileType(fileTypeRight)
	if err != nil {
		return nil, fmt.Errorf("right has unsupported file type: %w", err)
	}

	switch fileTypeRight {
	case dagger.FileTypeDirectory:
		directory := right.Directory(path)
		left = left.WithDirectory(path, directory)

	case dagger.FileTypeRegular:
		file := right.File(path)
		left = left.WithFile(path, file)

	}

	return left, nil
}

// Copy a specific entry from one directory to another, if the entry exists in both
func copyPathLeftExists(
	ctx context.Context,
	left *dagger.Directory,
	right *dagger.Directory,
	strategy MergeConflictStrategy,
	path string,
	fileTypeRight dagger.FileType,
) (*dagger.Directory, error) {
	var err error

	fileTypeLeft, err := left.Stat(path).FileType(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"could not get file type of path '%s': %w",
			path, err,
		)
	}

	err = assertAllowedFileType(fileTypeLeft)
	if err != nil {
		return nil, fmt.Errorf(
			"left path '%s' has unsupported file type: %w",
			path, err,
		)
	}

	err = assertAllowedFileType(fileTypeLeft)
	if err != nil {
		return nil, fmt.Errorf("left has unsupported file type: %w", err)
	}

	err = assertAllowedFileType(fileTypeRight)
	if err != nil {
		return nil, fmt.Errorf("right has unsupported file type: %w", err)
	}

	switch strategy {
	case ErrorOnConflict:
		return nil, fmt.Errorf("entry '%s' exists in both left and right", path)

	case KeepRight:
		switch fileTypeRight {
		case dagger.FileTypeDirectory:
			directory := right.Directory(path)
			return left.WithDirectory(path, directory), nil

		case dagger.FileTypeRegular:
			file := right.File(path)
			return left.WithFile(path, file), nil

		}

	case KeepLeft:
		return left, nil

	}

	return nil, fmt.Errorf("unexpected strategy '%s'", strategy)
}

func assertAllowedFileType(fileType dagger.FileType) error {
	switch fileType {
	case dagger.FileTypeDirectory:
		return nil
	case dagger.FileTypeRegular:
		return nil
	case dagger.FileTypeSymlink:
		return fmt.Errorf("symlink")
	case dagger.FileTypeUnknown:
		return fmt.Errorf("unknown")
	default:
		return fmt.Errorf("unexpected: %s", fileType)
	}
}
