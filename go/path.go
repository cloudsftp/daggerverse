package main

import "strings"

const (
	pathPrefix = "./"
)

func resolvePath(path string) string {
	if !strings.HasPrefix(path, pathPrefix) {
		path = pathPrefix + path
	}

	return path
}
