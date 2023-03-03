package collector

import (
	"github.com/mauricioabreu/video_samples/collector/watcher"
	"github.com/rs/zerolog/log"
)

func Collect(path string) {
	files, err := watcher.Watch(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize collector")
	}

	for file := range files {
		log.Debug().Msgf("File found: %s", file)
	}
}
