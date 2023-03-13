package collector

import (
	"github.com/google/uuid"

	"github.com/mauricioabreu/video_samples/collector/filesystem"
	"github.com/mauricioabreu/video_samples/collector/thumbnails"
	"github.com/mauricioabreu/video_samples/collector/watcher"
	"github.com/mauricioabreu/video_samples/config"
	"github.com/mauricioabreu/video_samples/store"
	"github.com/rs/zerolog/log"
)

func Collect(path string) {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get config")
	}
	rc := store.NewRedisCluster(cfg.RedisAddrs)
	files, err := watcher.Watch(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize collector")
	}

	for file := range files {
		log.Debug().Msgf("File found: %s", file)
		thumbnail, err := filesystem.NewFile(file)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read file: %s", file)
		}
		thumbnails.Insert(thumbnail, 60, uuid.NewString, rc)
	}
}
