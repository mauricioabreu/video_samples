package tasks

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/video_samples/extractor"
	"github.com/mauricioabreu/video_samples/extractor/inventory"
	"github.com/rs/zerolog/log"
)

const redisAddr = "127.0.0.1:6379"

func Enqueue(getStreamings func() ([]inventory.Streaming, error)) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	enqueueTasks(getStreamings, client)
	timer := time.NewTicker(30 * time.Second)
	for range timer.C {
		enqueueTasks(getStreamings, client)
	}
}

func enqueueTasks(getStreamings func() ([]inventory.Streaming, error), client *asynq.Client) {
	streamings, err := getStreamings()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get streamings")
		return
	}
	for _, stream := range streamings {
		task, err := NewExtractThumbsTask(extractor.ThumbOptions{
			Input:   stream.Playlist,
			Output:  "testvideo/thumbs/",
			Scale:   "-1:360",
			Quality: 5,
		})
		if err != nil {
			log.Error().Err(err).Msgf("Could not create task: %v", err)
			return
		}
		info, err := client.Enqueue(task)
		if err != nil {
			log.Error().Err(err).Msgf("Could not enqueue task: %v", err)
			return
		}
		log.Info().Msgf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}
}
