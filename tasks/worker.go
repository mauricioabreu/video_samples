package tasks

import (
	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/video_samples/config"
	"github.com/rs/zerolog/log"
)

const maxConcurrency = 10

func StartWorker(c *config.Config) {
	srv := asynq.NewServer(
		asynq.RedisClusterClientOpt{Addrs: c.RedisAddrs},
		asynq.Config{Concurrency: maxConcurrency},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(ThumbExtract, HandleThumbsExtractTask)
	mux.Use(loggingMiddleware)

	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msgf("Failed to start workers: %s", err)
	}
}
