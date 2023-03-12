package tasks

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const maxConcurrency = 10

func StartWorker() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: maxConcurrency},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(ThumbExtract, HandleThumbsExtractTask)
	mux.Use(loggingMiddleware)

	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msgf("Failed to start workers: %s", err)
	}
}
