package redis

import "github.com/redis/go-redis/v9"

func NewRedis(addrs []string) *redis.ClusterClient {
	rc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	return rc
}
