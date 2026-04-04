package main

import (
	"context"
	"fmt"
	"path/filepath"

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
	var err error

	for i, right := range rest {
		left, err = mergeTwoDirectories(ctx, left, right, strategy, "")
		if err != nil {
			return nil, fmt.Errorf("could not merge directory %d: %w", i+1, err)
		}
	}

	return left, nil
}

// Merge two directories into one
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
		return nil, fmt.Errorf("could not get entries right: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(currentPath, entry)

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

	err = assertAllowedFileType(fileTypeRight)
	if err != nil {
		return nil, fmt.Errorf(
			"right path '%s' has unsupported file type: %w",
			path, err,
		)
	}

	switch fileTypeRight {
	case dagger.FileTypeDirectory:
		switch fileTypeLeft {
		case dagger.FileTypeDirectory:
			return mergeTwoDirectories(ctx, left, right, strategy, path)

		case dagger.FileTypeRegular:
			switch strategy {
			case ErrorOnConflict:
				return nil, fmt.Errorf(
					"conflict for path '%s': right is directory and left is file",
					path,
				)

			case KeepRight:
				directory := right.Directory(path)
				return left.WithDirectory(path, directory), nil

			case KeepLeft:
				return left, nil

			default:
				return nil, fmt.Errorf("unexpected strategy: %s", strategy)
			}

		default:
			return nil, fmt.Errorf(
				"unexpected file type of path '%s' in left: %s",
				path, fileTypeLeft,
			)
		}

	case dagger.FileTypeRegular:
		switch fileTypeLeft {
		case dagger.FileTypeDirectory:
			switch strategy {
			case ErrorOnConflict:
				return nil, fmt.Errorf(
					"conflict for path '%s': right is file and left is directory",
					path,
				)

			case KeepRight:
				file := right.File(path)
				return left.WithFile(path, file), nil

			case KeepLeft:
				return left, nil
			}

		case dagger.FileTypeRegular:
			switch strategy {
			case ErrorOnConflict:
				return nil, fmt.Errorf(
					"conflict for path '%s': both files",
					path,
				)

			case KeepRight:
				file := right.File(path)
				return left.WithFile(path, file), nil

			case KeepLeft:
				return left, nil
			}

		default:
			return nil, fmt.Errorf(
				"unexpected file type of path '%s' in left: %s",
				path, fileTypeLeft,
			)
		}

	default:
		return nil, fmt.Errorf(
			"unexpected file type of path '%s' in right: %s",
			path, fileTypeRight,
		)

	}

	return nil, fmt.Errorf("unhandeled case at path '%s'", path)
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
