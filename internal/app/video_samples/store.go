package video_samples

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	log.Info(pong, err)
	return client
}

type store interface {
	GetThumb(streamName string) string
	GetThumbByTimestamp(streamName string, timestamp int64) string
	SaveThumb(stream Stream, timestamp int64, blob []byte) error
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore() redisStore {
	return redisStore{client: newClient()}
}

func (rs redisStore) GetThumb(streamName string) string {
	keys, err := rs.client.ZRevRange("thumbs/"+streamName, 0, 0).Result()
	if err != nil {
		log.Fatal(err)
	}
	thumb, err := rs.client.Get("thumbs/blob/" + keys[0]).Result()
	if err != nil {
		log.Fatal(err)
	}
	return thumb
}

func (rs redisStore) GetThumbByTimestamp(streamName string, timestamp int64) string {
	keys, err := rs.client.ZRangeByScore("thumbs/"+streamName, redis.ZRangeBy{
		Min:    strconv.FormatInt(timestamp, 10),
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		log.Fatal(err)
	}

	thumb, err := rs.client.Get("thumbs/blob/" + keys[0]).Result()
	if err != nil {
		log.Fatal(err)
	}
	return thumb
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
