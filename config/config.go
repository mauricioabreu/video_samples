package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RedisAddress  string `required:"true" split_words:"true"`
	RedisPassword string `required:"true" split_words:"true"`
	RedisDB       int    `default:"0" split_words:"true"`
}

func GetConfig() *Config {
	var c Config
	err := envconfig.Process("video_samples", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &c
}
