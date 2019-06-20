package thumbinator

import (
	"github.com/go-redis/redis"
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
}

type redisStore struct {
	client *redis.Client
}

func newRedisStore() redisStore {
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
