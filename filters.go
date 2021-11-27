package main

import (
	"os"
	"strings"
)

type filterer interface {
	filter(path string, fileInfo os.FileInfo) bool
}

type filterFunc func(path string, fileInfo os.FileInfo) bool

func (f filterFunc) filter(path string, fileInfo os.FileInfo) bool {
	return f(path, fileInfo)
}

func dotFilter(includeDot bool) filterer {
	return filterFunc(func(path string, fileInfo os.FileInfo) bool {
		return includeDot || !strings.HasPrefix(fileInfo.Name(), ".")
	})
}
