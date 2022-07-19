package video_samples

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	"github.com/mauricioabreu/video_samples/config"
	log "github.com/sirupsen/logrus"
)

func newClient(c *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddress,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	})
	pong, err := client.Ping().Result()
	log.Info(pong, err)
	return client
}

type store interface {
	GetThumb(streamName string) (string, error)
	GetThumbByTimestamp(streamName string, timestamp int64) (string, error)
	SaveThumb(stream Stream, timestamp int64, blob []byte) error
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(c *config.Config) redisStore {
	return redisStore{client: newClient(c)}
}

func (rs redisStore) GetThumb(streamName string) (string, error) {
	keys, err := rs.client.ZRevRange("thumbs/"+streamName, 0, 0).Result()
	if err != nil {
		log.Error(err)
		return "", err
	}
	thumb, err := rs.client.Get("thumbs/blob/" + keys[0]).Result()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return thumb, nil
}

func (rs redisStore) GetThumbByTimestamp(streamName string, timestamp int64) (string, error) {
	keys, err := rs.client.ZRangeByScore("thumbs/"+streamName, redis.ZRangeBy{
		Min:    strconv.FormatInt(timestamp, 10),
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		log.Error(err)
		return "", err
	}

	thumb, err := rs.client.Get("thumbs/blob/" + keys[0]).Result()
	if err != nil {
		log.Error(err)
		return "", err
	}
	return thumb, nil
}

func (rs redisStore) SaveThumb(stream Stream, timestamp int64, blob []byte) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	err = rs.client.ZAdd(fmt.Sprintf("thumbs/%s", stream.Name), redis.Z{Score: float64(timestamp), Member: id.String()}).Err()
	if err != nil {
		return err
	}

	if err = rs.client.Set(fmt.Sprintf("thumbs/blob/%s", id.String()), blob, time.Duration(stream.TTL)*time.Second).Err(); err != nil {
		return err
	}

	if err := rs.client.ZRemRangeByScore(fmt.Sprintf("thumbs/%s", stream.Name), "-inf", strconv.Itoa(int(timestamp)-stream.TTL)).Err(); err != nil {
		return err
	}

	return nil
}
