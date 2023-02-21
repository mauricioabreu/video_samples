package tasks

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func Start() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(ThumbExtract, HandleThumbsExtractTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msgf("Failed to start workers: %s", err)
	}
}
