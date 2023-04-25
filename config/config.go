package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	RedisAddr          string `env:"REDIS_ADDR,notEmpty"`
	EnqueueConcurrency int    `env:"WORKERS_CONCURRENCY" envDefault:"10"`
	InventoryAddress   string `env:"INVENTORY_ADDRESS,notEmpty"`
}

func GetConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to load configs: %w", err)
	}
	return cfg, nil
}
