package thumbnails

import (
	"context"
	"fmt"

	"github.com/mauricioabreu/video_samples/collector"
	"github.com/redis/go-redis/v9"
)

func Insert(file collector.File, uuid func() string, rc *redis.ClusterClient) error {
	uid := uuid()
	thumbsKey := fmt.Sprintf("thumbnails/%s", file.Dir)
	err := rc.ZAdd(context.TODO(), thumbsKey, redis.Z{Score: float64(file.ModTime), Member: uid}).Err()
	if err != nil {
		return fmt.Errorf("failed to insert into redis: %w", err)
	}

	return nil
}
