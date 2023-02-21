package main

import (
	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/video_samples/extractor"
	"github.com/mauricioabreu/video_samples/tasks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Print("Starting...")

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	task, err := tasks.NewExtractThumbsTask(extractor.ThumbOptions{
		Input:   "http://localhost:8080/big_buck_bunny/playlist.m3u8",
		Output:  "/tmp/thumbs/big_buck_bunny",
		Scale:   "-1:360",
		Quality: 5,
	})
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not enqueue task: %v", err)
	}
	log.Info().Msgf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
