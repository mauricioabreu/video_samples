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
