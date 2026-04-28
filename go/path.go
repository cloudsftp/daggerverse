package main

import "strings"

const (
	pathPrefix          = "./"
	recursivePathSuffix = "/..."
)

func resolvePath(path string) string {
	if !strings.HasPrefix(path, pathPrefix) {
		path = pathPrefix + path
	}

	return path
}

func resolveRecursivePath(path string) string {
	path = resolvePath(path)

	if !strings.HasSuffix(path, recursivePathSuffix) {
		path = path + recursivePathSuffix
	}

	return path
}
