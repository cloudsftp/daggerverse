package main

import (
	"context"
	"fmt"

	"dagger/merge-dirs/internal/dagger"
)

type MergeDirs struct{}

// Merge multiple directories into one
func (m *MergeDirs) MergeDirectories(
	ctx context.Context,
	dirs []*dagger.Directory,
) (*dagger.Directory, error) {
	if len(dirs) < 2 {
		return nil, fmt.Errorf("need at least 2 directories to merge")
	}

	first, rest := dirs[0], dirs[1:]
	var err error

	for i, next := range rest {
		first, err = mergeDirectories2(ctx, first, next)
		if err != nil {
			return nil, fmt.Errorf("could not merge directory %d: %w", i, err)
		}
	}

	return first, nil
}

// Merge two directories into one
func mergeDirectories2(
	ctx context.Context,
	first *dagger.Directory,
	second *dagger.Directory,
) (*dagger.Directory, error) {
	entries, err := second.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get entries from second directory: %w", err)
	}

	for _, path := range entries {
		first, err = copyPath(ctx, first, second, path)
		if err != nil {
			return nil, fmt.Errorf("could not copy at path '%s': %w", path, err)
		}
	}

	return first, nil
}

// Copy a specific path from one directory to another
func copyPath(
	ctx context.Context,
	target *dagger.Directory,
	source *dagger.Directory,
	path string,
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
