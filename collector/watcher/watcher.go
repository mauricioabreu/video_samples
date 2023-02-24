package watcher

import (
	"fmt"

	"github.com/rjeczalik/notify"
	"github.com/rs/zerolog/log"
)

func Watch(path string) error {
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(path, c, WriteEvent); err != nil {
		return fmt.Errorf("failed to watch %s: %w", path, err)
	}
	defer notify.Stop(c)

	for {
		e := <-c
		log.Info().Msgf("Event %v received for: %s", e.Event().String(), e.Path())
	}
}
