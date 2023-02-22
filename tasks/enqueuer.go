package tasks

import (
	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/video_samples/extractor"
	"github.com/rs/zerolog/log"
)

const redisAddr = "127.0.0.1:6379"

func Enqueue() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	task, err := NewExtractThumbsTask(extractor.ThumbOptions{
		Input:   "http://localhost:8080/output.m3u8",
		Output:  "testvideo/thumbs/",
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
