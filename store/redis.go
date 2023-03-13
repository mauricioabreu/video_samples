package store

import "github.com/redis/go-redis/v9"

func NewRedisCluster(addrs []string) *redis.ClusterClient {
	rc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	return rc
}
