package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mauricioabreu/video_samples/extractor"
	"github.com/rs/zerolog/log"
)

const (
	ThumbExtract = "thumbs:extract"
)

func NewExtractThumbsTask(to extractor.ThumbOptions) (*asynq.Task, error) {
	payload, err := json.Marshal(to)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(ThumbExtract, payload), nil
}

func HandleThumbsExtractTask(ctx context.Context, t *asynq.Task) error {
	var opts extractor.ThumbOptions
	if err := json.Unmarshal(t.Payload(), &opts); err != nil {
		return fmt.Errorf("failed to process payload: %w %w", err, asynq.SkipRetry)
	}
	log.Info().Msgf("Extracting thumbs from video URL: %s", opts.Input)
	if err := extractor.ExtractThumbs("colors", opts, extractor.RunCmd); err != nil {
		return fmt.Errorf("failed to run the extractor: %w %w", err, asynq.SkipRetry)
	}
	return nil
}
