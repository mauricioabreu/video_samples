package tasks

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		log.Debug().Msgf("Started processing %s", t.Type())
		if err := h.ProcessTask(ctx, t); err != nil {
			log.Error().Err(err).Msgf("Failed to run task: %q", t.Type())
			return err
		}
		log.Debug().Msgf("Finished processing %s: Elapsed time %v", t.Type(), time.Since(start))
		return nil
	})
}
