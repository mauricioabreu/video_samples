package tasks

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/video_samples/config"
	"github.com/mauricioabreu/video_samples/extractor"
	"github.com/mauricioabreu/video_samples/extractor/inventory"
	"github.com/rs/zerolog/log"
)

const runEvery = 30 * time.Second

func Enqueue(getStreamings func() ([]inventory.Streaming, error)) {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	client := asynq.NewClient(asynq.RedisClusterClientOpt{Addrs: cfg.RedisAddrs})
	defer client.Close()

	enqueueTasks(getStreamings, client)
	timer := time.NewTicker(runEvery)
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
		info, err := client.Enqueue(task, asynq.TaskID(generateID("thumb", stream.Name)))
		if err != nil {
			log.Error().Err(err).Msgf("Could not enqueue task: %v", err)
			return
		}
		log.Info().Msgf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}
}

func generateID(feature, name string) string {
	return fmt.Sprintf("%s_%s", feature, name)
}
