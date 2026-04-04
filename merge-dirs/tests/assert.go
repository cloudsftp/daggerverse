package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"
)

func assertDirectory(
	ctx context.Context,
	expected *dagger.Directory,
	actual *dagger.Directory,
) error {
	expectedFiles, err := collectAllFiles(ctx, expected)
	if err != nil {
		return fmt.Errorf("could not get file names of expected directory: %w", err)
	}

	actualFiles, err := collectAllFiles(ctx, actual)
	if err != nil {
		return fmt.Errorf("could not get file names of actual directory: %w", err)
	}

	for path, expectedContent := range expectedFiles {
		actualContent, ok := actualFiles[path]
		if !ok {
			return fmt.Errorf("file at path '%s' does not exist", path)
		}

		if expectedContent != actualContent {
			return fmt.Errorf(
				"unexpected content in file at path '%s': expected '%s', got '%s'",
				path, expectedContent, actualContent,
			)
		}
	}

	for path := range actualFiles {
		_, ok := expectedFiles[path]
		if !ok {
			return fmt.Errorf("did not expect file at '%s'", path)
		}
	}

	return nil
}
