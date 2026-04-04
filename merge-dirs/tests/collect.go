package main

import (
	"context"
	"fmt"
	pathlib "path"

	"dagger/tests/internal/dagger"
)

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
		path := pathlib.Join(currentPath, entry)

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
