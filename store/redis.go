package store

import "github.com/redis/go-redis/v9"

func NewRedis(addr string) *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return rc
}
