package main

import (
	"fmt"

	"dagger/tests/internal/dagger"
)

func buildDirectory(files map[string]any) *dagger.Directory {
	dir := dag.Directory()

	for name, value := range files {
		switch v := value.(type) {
		case string:
			dir = dir.WithNewFile(name, v)
		case map[string]any:
			dir = dir.WithDirectory(name, buildDirectory(v))
		default:
			panic(fmt.Sprintf("unexpected type for %s: %T", name, v))
		}
	}

	return dir
}
