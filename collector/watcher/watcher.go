package watcher

import (
	"fmt"
	"path/filepath"

	"github.com/rjeczalik/notify"
	"github.com/samber/lo"
)

type File struct {
	Path    string
	Dir     string
	ModTime int64
}

// Magic number. notify demands for a buffered channel because it does not block sending
// to the channel
var filesBuffer = 200

// Watch files in a given path and sends events to channels to be
// processed later
func Watch(path string) (<-chan string, error) {
	files := make(chan string, filesBuffer)

	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(path, c, WriteEvent); err != nil {
		return nil, fmt.Errorf("failed to watch %s: %w", path, err)
	}

	go func() {
		for {
			e := <-c
			path := e.Path()
			if MatchImage(path) {
				files <- path
			}
		}
	}()

	return files, nil
}

// MatchExt checks if a given path matches a list of patterns
func MatchExt(path string, patterns []string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	return lo.Contains(patterns, ext[1:])
}

func MatchImage(path string) bool {
	return MatchExt(path, []string{"jpg", "png"})
}
