package watcher

import (
	"fmt"
	"path/filepath"

	"github.com/rjeczalik/notify"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

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
	defer notify.Stop(c)

	go func() {
		for {
			e := <-c
			log.Info().Msgf("Event %v received for: %s", e.Event().String(), e.Path())
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
