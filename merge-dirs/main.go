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
	// +default="KEEP_LEFT"
	// Conflict resolution strategy: left (default), right, or error
	strategy MergeConflictStrategy,
) (*dagger.Directory, error) {
	if len(dirs) < 2 {
		return nil, fmt.Errorf("need at least 2 directories to merge, got %d", len(dirs))
	}

	first, rest := dirs[0], dirs[1:]
	var err error

	for i, next := range rest {
		first, err = mergeDirectories2(ctx, first, next, strategy)
		if err != nil {
			return nil, fmt.Errorf("could not merge directory %d: %w", i, err)
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
) (*dagger.Directory, error) {
	entries, err := right.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get entries from second directory: %w", err)
	}

	for _, path := range entries {
		left, err = copyPath(ctx, left, right, path, strategy)
		if err != nil {
			return nil, fmt.Errorf("could not copy at path '%s': %w", path, err)
		}
	}

	return left, nil
}

// Copy a specific path from one directory to another
func copyPath(
	ctx context.Context,
	target *dagger.Directory,
	source *dagger.Directory,
	path string,
	strategy MergeConflictStrategy,
) (*dagger.Directory, error) {
	fileType, err := source.Stat(path).FileType(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get file type of path: %w", err)
	}

	switch fileType {
	case dagger.FileTypeDirectory:
		directory := source.Directory(path)
		target = target.WithDirectory(path, directory)

	case dagger.FileTypeRegular:
		file := source.File(path)
		target = target.WithFile(path, file)

	case dagger.FileTypeSymlink:
		return nil, fmt.Errorf("symlinks are not supported")

	case dagger.FileTypeUnknown:
		return nil, fmt.Errorf("unknown file type")

	default:
		return nil, fmt.Errorf("unexpected file type: %#v", fileType)
	}

	return target, nil
}
