package thumbnails

import (
	"context"
	"fmt"
	"time"

	"github.com/mauricioabreu/video_samples/collector"
	"github.com/redis/go-redis/v9"
)

func Insert(file collector.File, uuid func() string, rc *redis.ClusterClient) error {
	thumbId := fmt.Sprintf("blob/%s", uuid())
	thumbsKey := fmt.Sprintf("thumbnails/%s", file.Dir)

	err := rc.ZAdd(context.TODO(), thumbsKey, redis.Z{Score: float64(file.ModTime), Member: thumbId}).Err()
	if err != nil {
		return fmt.Errorf("failed to insert into redis: %w", err)
	}

	expiration := time.Duration(10) * time.Second
	if err := rc.Set(context.TODO(), thumbId, "data", expiration).Err(); err != nil {
		return fmt.Errorf("failed to insert into redis: %w", err)
	}

	return nil
}
