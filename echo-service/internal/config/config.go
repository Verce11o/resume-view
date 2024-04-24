package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env                 string `env:"env"`
	Server              Server
	ViewServiceEndpoint string `env:"VIEW_SERVICE_ENDPOINT" env-default:"localhost:3007"`
	ClientTimeout       string `env:"CLIENT_TIMEOUT" env-default:"5s"`
	RetriesCount        string `env:"RETRIES_COUNT" env-default:"3"`
}

type Server struct {
	Port string `env:"SERVER_PORT" env-default:"3008"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return &cfg
}
