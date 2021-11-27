package main

import (
	"os"
	"strings"
)

type filterer interface {
	filter(path string, entry os.DirEntry) bool
}

type filterFunc func(path string, entry os.DirEntry) bool

func (f filterFunc) filter(path string, entry os.DirEntry) bool {
	return f(path, entry)
}

func dotFilter(includeDot bool) filterer {
	return filterFunc(func(path string, entry os.DirEntry) bool {
		return includeDot || !strings.HasPrefix(entry.Name(), ".")
	})
}
