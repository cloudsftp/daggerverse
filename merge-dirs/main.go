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
		first, err = mergeDirectories2(ctx, first, next, strategy)
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
) (*dagger.Directory, error) {
	entries, err := right.Entries(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get entries from second directory: %w", err)
	}

	for _, path := range entries {
		left, err = copyEntry(ctx, left, right, path, strategy)
		if err != nil {
			return nil, fmt.Errorf("could not copy at path '%s': %w", path, err)
		}
	}

	return left, nil
}

// Copy a specific path from one directory to another
func copyEntry(
	ctx context.Context,
	left *dagger.Directory,
	right *dagger.Directory,
	entry string,
	strategy MergeConflictStrategy,
) (*dagger.Directory, error) {
	fileTypeLeft, err := left.Stat(entry).FileType(ctx)
	existsLeft := err == nil

	fileTypeRight, err := right.Stat(entry).FileType(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get file type of path: %w", err)
	}

	if existsLeft {
		return copyEntryLeftExists(ctx, left, fileTypeLeft, right, fileTypeRight, entry, strategy)
	}

	err = assertAllowedFileType(fileTypeRight)
	if err != nil {
		return nil, fmt.Errorf("right has unsupported file type: %w", err)
	}

	switch fileTypeRight {
	case dagger.FileTypeDirectory:
		directory := right.Directory(entry)
		left = left.WithDirectory(entry, directory)

	case dagger.FileTypeRegular:
		file := right.File(entry)
		left = left.WithFile(entry, file)

	}

	return left, nil
}

// Copy a specific entry from one directory to another, if the entry exists in both
func copyEntryLeftExists(
	ctx context.Context,
	left *dagger.Directory,
	fileTypeLeft dagger.FileType,
	right *dagger.Directory,
	fileTypeRight dagger.FileType,
	entry string,
	strategy MergeConflictStrategy,
) (*dagger.Directory, error) {
	var err error

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
		return nil, fmt.Errorf("entry '%s' exists in both left and right", entry)

	case KeepRight:
		switch fileTypeRight {
		case dagger.FileTypeDirectory:
			directory := right.Directory(entry)
			return left.WithDirectory(entry, directory), nil

		case dagger.FileTypeRegular:
			file := right.File(entry)
			return left.WithFile(entry, file), nil

		}

	case KeepLeft:
		return nil, fmt.Errorf("merge strategy keep left not yet implemented")
	}

	return nil, nil
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
