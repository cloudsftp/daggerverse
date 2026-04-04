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
) (map[string]string, error) {
	return recursivelyCollectFiles(ctx, directory, "", map[string]string{})
}

func recursivelyCollectFiles(
	ctx context.Context,
	directory *dagger.Directory,
	currentPath string,
	files map[string]string,
) (map[string]string, error) {
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
			filesInDirectory, err := recursivelyCollectFiles(ctx, directory, path, files)
			if err != nil {
				return nil, err
			}
			for path, content := range filesInDirectory {
				files[path] = content
			}

		case dagger.FileTypeRegular:
			content, err := directory.File(path).Contents(ctx)
			if err != nil {
				return nil, fmt.Errorf(
					"could not get content of file at path '%s': %w",
					path, err,
				)
			}
			files[path] = content

		case dagger.FileTypeSymlink:
			return nil, fmt.Errorf("path '%s' is a symlink", path)

		case dagger.FileTypeUnknown:
			return nil, fmt.Errorf("unknown file type of path '%s'", path)

		default:
			return nil, fmt.Errorf("file type of path '%s' unexpected", path)
		}
	}

	return files, nil
}
