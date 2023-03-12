package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	RedisAddrs         []string `env:"REDIS_ADDRS,notEmpty" envSeparator:","`
	EnqueueConcurrency int      `env:"WORKERS_CONCURRENCY" envDefault:"10"`
	InventoryAddress   string   `env:"INVENTORY_ADDRESS,notEmpty"`
}

func GetConfig() (*Config, error) {
	var cfg *Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load configs: %w", err)
	}
	return cfg, nil
}
